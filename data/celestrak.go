package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Funkit/tle-provider/apierror"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/Funkit/tle-provider/utils"
)

// CelestrakData structure used to parse Celestrak response when using the JSON API
type CelestrakData struct {
	ObjectName         string  `json:"OBJECT_NAME"`
	ObjectID           string  `json:"OBJECT_ID"`
	Epoch              string  `json:"EPOCH"`
	MeanMotion         float64 `json:"MEAN_MOTION"`
	Eccentricity       float64 `json:"ECCENTRICITY"`
	Inclination        float64 `json:"INCLINATION"`
	RaOfASCMode        float64 `json:"RA_OF_ASC_NODE"`
	ArgOfPericenter    float64 `json:"ARG_OF_PERICENTER"`
	MeanAnomaly        float64 `json:"MEAN_ANOMALY"`
	EphemerisType      int     `json:"EPHEMERIS_TYPE"`
	ClassificationType string  `json:"CLASSIFICATION_TYPE"`
	NORADCatID         int     `json:"NORAD_CAT_ID"`
	ElementSetNo       int     `json:"ELEMENT_SET_NO"`
	RevAtEpoch         int     `json:"REV_AT_EPOCH"`
	BStar              float64 `json:"BSTAR"`
	MeanMotionDOT      float64 `json:"MEAN_MOTION_DOT"`
	MeanMotionDDOT     float64 `json:"MEAN_MOTION_DDOT"`
}

// CelestrakClient implementation of the Source interface for Celestrak
type CelestrakClient struct {
	httpClient         *http.Client
	AllSatellitesURL   string
	GeoSatellitesURL   string
	OrbitalData        []CelestrakData
	TwoLineElements    []*Satellite
	TwoLineElementsMap map[string]*Satellite
	LastCelestrakPull  time.Time
	UpdatePeriod       float64
	mu                 sync.RWMutex
}

// NewCelestrakClient Generates a new CelestrakClient from the information in the configuration file
func NewCelestrakClient(allSatellitesURL, geoSatellitesURL string, refreshRateHours int) *CelestrakClient {

	return &CelestrakClient{
		httpClient:        &http.Client{},
		AllSatellitesURL:  allSatellitesURL,
		GeoSatellitesURL:  geoSatellitesURL,
		OrbitalData:       []CelestrakData{},
		LastCelestrakPull: time.Date(1970, 01, 01, 0, 0, 0, 1, time.UTC),
		UpdatePeriod:      float64(refreshRateHours),
	}
}

//GetDataSource return server data source
func (cc *CelestrakClient) GetDataSource() string {
	return "celestrak"
}

//GetConfig return server configuration
func (cc *CelestrakClient) GetConfig() (map[string]interface{}, error) {
	return nil, nil
}

// getCelestrakData Get data from celestrak
func (cc *CelestrakClient) getCelestrakData() ([]CelestrakData, error) {
	req, err := http.NewRequest(http.MethodGet, cc.AllSatellitesURL, nil)
	if err != nil {
		return nil, err
	}

	response, err := cc.httpClient.Do(req)
	if err != nil {
		return nil, apierror.Wrap(err, apierror.ErrInternal)
	}
	defer response.Body.Close()

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, apierror.Wrap(err, apierror.ErrInternal)
	}

	if response.StatusCode >= 400 {
		return nil, apierror.Wrap(fmt.Errorf("failed to query data from celestrak, response error code = %v", response.StatusCode), apierror.ErrNotFound)
	}

	var output []CelestrakData

	if err := json.Unmarshal(respBody, &output); err != nil {
		return nil, apierror.Wrap(err, apierror.ErrInternal)
	}

	return output, nil
}

func (cc *CelestrakClient) Update(done <-chan struct{}, period time.Duration) {
	if err := cc.update(); err != nil {
		log.Println(err.Error())
	}
	go func() {
		for {
			select {
			case <-done:
				break
			case <-time.After(period):
				if err := cc.update(); err != nil {
					log.Println(err.Error())
				}
			}
		}
	}()
}

// GetData Implementation of the Source interface for Celestrak
func (cc *CelestrakClient) GetData() ([]*Satellite, error) {
	cc.mu.RLock()
	defer cc.mu.RUnlock()
	if len(cc.TwoLineElements) == 0 {
		return nil, apierror.Wrap(fmt.Errorf("No satellite found"), apierror.ErrNotFound)
	}
	return cc.TwoLineElements, nil
}

func (cc *CelestrakClient) update() error {
	satData, err := cc.getCelestrakData()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	output := make(chan *Satellite, len(satData))

	for i := 0; i < len(satData); i++ {
		wg.Add(1)
		go func(data CelestrakData) {
			defer wg.Done()
			sat, _ := convertToTLE(data)
			output <- &sat
		}(satData[i])
	}

	wg.Wait()
	close(output)
	cc.mu.Lock()
	defer cc.mu.Unlock()
	var tleList []*Satellite
	cc.TwoLineElementsMap = make(map[string]*Satellite)
	for element := range output {
		tleList = append(tleList, element)
		cc.TwoLineElementsMap[element.SatelliteName] = element

	}
	cc.TwoLineElements = tleList
	cc.LastCelestrakPull = time.Now()
	log.Printf("data successfully pulled from celestrak at %s\n", time.Now().Format("2006-01-02T15:04:05Z"))

	return nil
}

func (cc *CelestrakClient) GetSatellite(satelliteName string) chan *SatelliteErr {
	output := make(chan *SatelliteErr)
	go func() {
		cc.mu.RLock()
		defer cc.mu.RUnlock()
		if cc.TwoLineElementsMap[satelliteName] == nil {
			output <- &SatelliteErr{
				Err: apierror.Wrap(fmt.Errorf("Satellite %v not found", satelliteName), apierror.ErrNotFound),
				Sat: nil,
			}
		} else {
			output <- &SatelliteErr{
				Err: nil,
				Sat: cc.TwoLineElementsMap[satelliteName],
			}
		}
	}()
	return output
}

// convertToTLE converts JSON from Celestrak GP prototype to TLE. See wikipedia page for two line elements for an explanation of the fields
func convertToTLE(data CelestrakData) (Satellite, error) {

	cosparID, err := objectIDToCOSPARID(data.ObjectID)
	if err != nil {
		return Satellite{}, apierror.Wrap(err, apierror.ErrRender)
	}

	year2Digits, err := getLast2DigitsOfYear(data.Epoch)
	if err != nil {
		return Satellite{}, apierror.Wrap(err, apierror.ErrRender)
	}

	dayOfYear, err := getDayOfYear(data.Epoch)
	if err != nil {
		return Satellite{}, apierror.Wrap(err, apierror.ErrRender)
	}

	// Line 1 formatting

	line1Items := []string{
		"1 ", // field 1
		fmt.Sprintf("%05d", data.NORADCatID) + data.ClassificationType, // field 2, 3
		" ",
		cosparID, // field 4, 5, 6
		" ",
		year2Digits + dayOfYear, // field 7, 8
		" ",
		formatMeanMotionDOT(data.MeanMotionDOT), // field 9
		" ",
		formatWithoutDecimalPoint(data.MeanMotionDDOT), // field 10
		" ",
		formatWithoutDecimalPoint(data.BStar), // field 11
		" 0 ",                                 // field 12
		formatWithLeadingSpaces(data.ElementSetNo), // field 13
	}

	tleLine1WithoutChecksum := strings.Join(line1Items, "")

	checksum1, err := checksumAsString(tleLine1WithoutChecksum)
	if err != nil {
		return Satellite{}, err
	}
	tleLine1 := tleLine1WithoutChecksum + checksum1

	// Line 2 formatting

	line2Items := []string{
		"2 ",                                 // field 1
		fmt.Sprintf("%05d", data.NORADCatID), // field 2
		" ",
		formatAngles(data.Inclination), // field 3
		" ",
		formatAngles(data.RaOfASCMode), // field 4
		" ",
		formatEccentricity(data.Eccentricity), // field 5
		" ",
		formatAngles(data.ArgOfPericenter), // field 6
		" ",
		formatAngles(data.MeanAnomaly), // field 7
		" ",
		formatMeanMotion(data.MeanMotion), // field 8
		formatRevNumber(data.RevAtEpoch),  // field 9
	}

	tleLine2WithoutChecksum := strings.Join(line2Items, "")

	checksum2, err := checksumAsString(tleLine2WithoutChecksum)
	if err != nil {
		return Satellite{}, err
	}
	tleLine2 := tleLine2WithoutChecksum + checksum2

	return Satellite{
		SatelliteName: data.ObjectName,
		NORADID:       int(data.NORADCatID),
		TLELine1:      tleLine1,
		TLELine2:      tleLine2,
	}, nil
}

func objectIDToCOSPARID(objectID string) (string, error) {
	re := regexp.MustCompile(`[0-9]{2}(.+)-(.+)`)
	matchResults := re.FindAllSubmatch([]byte(objectID), -1)
	if (len(matchResults) != 1) || (len(matchResults[0]) != 3) {
		return "", fmt.Errorf("Could not convert %s to COSPAR ID", objectID)
	}

	// Add trailing spaces
	trailingSpaceNumber := 6 - len(matchResults[0][2])
	trailingSpaces := ""
	for i := 0; i < trailingSpaceNumber; i++ {
		trailingSpaces = trailingSpaces + " "
	}

	return string(matchResults[0][1]) + string(matchResults[0][2]) + trailingSpaces, nil
}

func getLast2DigitsOfYear(epoch string) (string, error) {
	re := regexp.MustCompile(`[0-9]{2}([0-9]{2}).+`)
	matchResults := re.FindAllSubmatch([]byte(epoch), -1)
	if (len(matchResults) != 1) || (len(matchResults[0]) != 2) {
		return "", fmt.Errorf("cannot parse last 2 digits of year")
	}

	return string(matchResults[0][1]), nil
}

func getDayOfYear(epoch string) (string, error) {
	layout := "2006-01-02T15:04:05.000000"
	t, err := time.Parse(layout, epoch)
	if err != nil {
		return "", err
	}

	dayOfYearFloat := float64(t.YearDay()) + float64(t.Hour())/24 + float64(t.Minute())/(60*24) + float64(t.Second())/(3600*24) + float64(t.Nanosecond())/(1000000000*3600*24)

	// For some reason the fmt formatting %03.8f does not add leading zeroes, so manual padding for single and double digits added manually
	dayString := fmt.Sprintf("%.8f", dayOfYearFloat)
	if dayOfYearFloat < 100 {
		dayString = "0" + dayString
	}
	if dayOfYearFloat < 10 {
		dayString = "0" + dayString
	}
	return dayString, nil
}

func checksumAsString(line string) (string, error) {
	var checksum int
	if len(line) != 68 {
		errorText := fmt.Sprintf("Generated line is not 68 characters long. Line: %s", line)
		return "", errors.New(errorText)
	}
	for i := 0; i < 68; i++ {
		if line[i] == '-' {
			checksum++
		} else if line[i] != ' ' && line[i] != '.' && line[i] != '+' && !utils.IsALetter(line[i]) {
			checksum = checksum + int(line[i]) - 48
		}
	}
	return fmt.Sprintf("%d", checksum%10), nil
}
