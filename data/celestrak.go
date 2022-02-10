package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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
	httpClient        *http.Client
	AllSatellitesURL  string
	GeoSatellitesURL  string
	OrbitalData       []CelestrakData
	LastCelestrakPull time.Time
	UpdatePeriod      float64
	Lock              sync.Mutex
}

// NewCelestrakClient Generates a new CelestrakClient from the information in the configuration file
func NewCelestrakClient(config *utils.Info) *CelestrakClient {
	return &CelestrakClient{
		httpClient:        &http.Client{},
		AllSatellitesURL:  config.CelestrakURLs.AllSatellites,
		GeoSatellitesURL:  config.CelestrakURLs.GeoSatellites,
		OrbitalData:       []CelestrakData{},
		LastCelestrakPull: time.Date(1970, 01, 01, 0, 0, 0, 1, time.UTC),
		UpdatePeriod:      1,
	}
}

// GetData Implementation of the Source interface for Celestrak
func (cc *CelestrakClient) GetData() ([]Satellite, error) {

	elapsedTime := time.Now().Sub(cc.LastCelestrakPull)

	fmt.Println("Elapsed time: ", elapsedTime)
	fmt.Println("Last Celestrak Pull: ", cc.LastCelestrakPull)

	if elapsedTime.Hours() >= cc.UpdatePeriod {
		fmt.Println("Pulling from Celestrak")
		cc.LastCelestrakPull = time.Now()
		if err := cc.GetCelestrakData(); err != nil {
			return []Satellite{}, err
		}
	}

	var tleList []Satellite

	for _, element := range cc.OrbitalData {
		sat, err := convertToTLE(element)
		if err != nil {
			return []Satellite{}, err
		}
		tleList = append(tleList, sat)
	}

	return tleList, nil
}

// GetCelestrakData Get data from celestrak
func (cc *CelestrakClient) GetCelestrakData() error {
	req, err := http.NewRequest(http.MethodGet, cc.AllSatellitesURL, nil)
	if err != nil {
		return err
	}

	response, err := cc.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode >= 400 {
		return errors.New(fmt.Sprintln("HTTP code:", response.StatusCode, "response body:", string(respBody)))
	}

	cc.OrbitalData = nil
	return json.Unmarshal(respBody, &cc.OrbitalData)
}

// convertToTLE converts JSON from Celestrak GP prototype to TLE. See wikipedia page for two line elements for an explanation of the fields
func convertToTLE(data CelestrakData) (Satellite, error) {

	cosparID, err := objectIDToCOSPARID(data.ObjectID)
	if err != nil {
		return Satellite{}, err
	}

	year2Digits, err := getLast2DigitsOfYear(data.Epoch)
	if err != nil {
		return Satellite{}, err
	}

	dayOfYear, err := getDayOfYear(data.Epoch)
	if err != nil {
		return Satellite{}, err
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

	tleLine1 := tleLine1WithoutChecksum + checksumAsString(tleLine1WithoutChecksum)

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

	tleLine2 := tleLine2WithoutChecksum + checksumAsString(tleLine2WithoutChecksum)

	return Satellite{
		SatelliteName: data.ObjectName,
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

	return fmt.Sprintf("%3.8f", dayOfYearFloat), nil
}

func checksumAsString(line string) string {
	var checksum int
	for i := 0; i < 68; i++ {
		if line[i] == '-' {
			checksum++
		} else if line[i] != ' ' && line[i] != '.' && line[i] != '+' && !utils.IsALetter(line[i]) {
			checksum = checksum + int(line[i]) - 48
		}
	}
	return fmt.Sprintf("%d", checksum%10)
}