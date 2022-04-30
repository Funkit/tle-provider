package utils

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

//Info Receiving structure when parsing the configuration file
type Info struct {
	ServerPort             int                    `yaml:"server_port"`
	DataSource             string                 `yaml:"data_source"`
	CelestrakConfiguration CelestrakConfiguration `yaml:"celestrak_configuration"`
}

type CelestrakConfiguration struct {
	AllSatellitesURL string `yaml:"all_satellites_url"`
	GeoSatellitesURL string `yaml:"geo_satellites_url"`
	RefreshRateHours int    `yaml:"celestrak_refresh_rate_hours"`
}

// IsALetter Check if char is an ASCII letter or not
func IsALetter(element byte) bool {
	val := ((element >= 'a') && (element <= 'z')) || ((element >= 'A') && (element <= 'Z'))

	return val
}

//GenericJSONParsing Parse into struct
func GenericJSONParsing[T any](filePath string) (T, error) {

	var x T

	rawContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return x, fmt.Errorf("Error when reading the configuration file: %w", err)
	}

	err = json.Unmarshal(rawContent, &x)
	if err != nil {
		return x, fmt.Errorf("Error when unmarshalling the JSON file: %w", err)
	}

	return x, err
}

//GenericYAMLParsing Parse into struct
func GenericYAMLParsing[T any](filePath string) (T, error) {

	var x T

	rawContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return x, fmt.Errorf("Error when reading the configuration file: %w", err)
	}

	err = yaml.Unmarshal(rawContent, &x)
	if err != nil {
		return x, fmt.Errorf("Error when unmarshalling the YAML file: %w", err)
	}

	return x, err
}

//ValueEqual check that 2 pointers to comparable items point to values that are equal
func ValueEqual[T comparable](item1, item2 *T) bool {
	if item1 == nil {
		if item2 != nil {
			return false
		}
	} else {
		if item2 == nil {
			return false
		} else {
			if *item1 != *item2 {
				return false
			}
		}
	}
	return true
}
