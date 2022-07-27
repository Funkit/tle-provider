package data

import (
	"fmt"
	"github.com/Funkit/go-utils/utils"
	"github.com/Funkit/tle-provider/apierror"
	"log"
	"strconv"
	"sync"
	"time"
)

type FileSource struct {
	filePath           string
	TwoLineElements    []Satellite
	TwoLineElementsMap map[string]Satellite
	LastCelestrakPull  time.Time
	UpdatePeriod       float64
	mu                 sync.RWMutex
}

func NewFileSource(filePath string, refreshRateSeconds int) *FileSource {
	return &FileSource{
		filePath:          filePath,
		LastCelestrakPull: time.Date(1970, 01, 01, 0, 0, 0, 1, time.UTC),
		UpdatePeriod:      float64(refreshRateSeconds),
	}
}

func (fs *FileSource) Update(done <-chan struct{}, period time.Duration) {
	if err := fs.update(); err != nil {
		log.Println(err.Error())
	}
	go func() {
		for {
			select {
			case <-done:
				break
			case <-time.After(period):
				if err := fs.update(); err != nil {
					log.Println(err.Error())
				}
			}
		}
	}()
}

func (fs *FileSource) update() error {

	tleList, err := fs.extractSatelliteData()
	if err != nil {
		return err
	}

	fs.mu.Lock()
	defer fs.mu.Unlock()
	fs.TwoLineElements = tleList
	return nil
}

func (fs *FileSource) GetData() ([]Satellite, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	if len(fs.TwoLineElements) == 0 {
		return nil, apierror.Wrap(fmt.Errorf("No satellite found"), apierror.ErrNotFound)
	}
	return fs.TwoLineElements, nil
}

func (fs *FileSource) GetSatellite(satelliteName string) chan SatelliteErr {
	output := make(chan SatelliteErr)
	go func() {
		fs.mu.RLock()
		defer fs.mu.RUnlock()
		if fs.TwoLineElementsMap[satelliteName].IsNull() {
			output <- SatelliteErr{
				Err: apierror.Wrap(fmt.Errorf("Satellite %v not found", satelliteName), apierror.ErrNotFound),
				Sat: Satellite{},
			}
		} else {
			output <- SatelliteErr{
				Err: nil,
				Sat: fs.TwoLineElementsMap[satelliteName],
			}
		}
	}()
	return output
}

func (fs *FileSource) GetDataSource() string {
	return "file"
}

func (fs *FileSource) GetConfig() (map[string]interface{}, error) {
	return nil, nil
}

func (fs *FileSource) extractSatelliteData() ([]Satellite, error) {
	fileLines, err := utils.GetFileAsLines(fs.filePath)
	if err != nil {
		return nil, err
	}

	var output []Satellite

	i := 0
	for i+2 < len(fileLines) {
		if len(fileLines[i+1]) != 69 {
			return nil, fmt.Errorf("TLE line 1 for item %v has wrong format, expected 69 characters, got %v", i, len(fileLines[i+1]))
		}
		if len(fileLines[i+2]) != 69 {
			return nil, fmt.Errorf("TLE line 2 for item %v has wrong format, expected 69 characters, got %v", i, len(fileLines[i+2]))
		}

		noradID, err := strconv.Atoi(fileLines[i+2][2:7])
		if err != nil {
			return nil, fmt.Errorf("error parsing the NORAD ID from the first TLE line: %v could not be cast as an int", fileLines[i+1][2:6])
		}

		output = append(output, Satellite{
			SatelliteName: removeTrailingSpaces(fileLines[i]),
			NORADID:       noradID,
			TLELine1:      fileLines[i+1],
			TLELine2:      fileLines[i+2],
		})
		i = i + 3
	}

	return output, nil
}

func removeTrailingSpaces(source string) string {

	j := len(source) - 1

	for j >= 0 {
		if source[j] != ' ' {
			break
		}
		j--
	}

	if j == 0 {
		return source
	} else if j < 0 {
		return ""
	} else {
		return source[0 : j+1]
	}
}
