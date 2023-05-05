package data

import (
	"fmt"
	"strconv"

	"github.com/Funkit/go-utils/utils"
)

type FileSource struct {
	filePath string
}

func NewFileSource(filePath string) *FileSource {
	return &FileSource{
		filePath: filePath,
	}
}

func (fs *FileSource) GetData() ([]Satellite, error) {
	tleList, err := fs.extractSatelliteData()
	if err != nil {
		return nil, err
	}
	return tleList, nil
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
