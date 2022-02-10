package utils

import (
	"fmt"
	"io/ioutil"

	"github.com/alexflint/go-arg"
	"gopkg.in/yaml.v2"
)

var args struct {
	ConfigFilePath string `arg:"required,-c" help:"path to the configuration file"`
}

//Info server configuration
type Info struct {
	CelestrakURLs             CelestrakURLs `yaml:"celestrak_urls"`
	CelestrakRefreshRateHours string        `yaml:"celestrak_refresh_rate_hours"`
	ServerPort                int           `yaml:"server_port"`
	DataSource                string        `yaml:"data_source"`
}

//CelestrakURLs Celestrak addresses
type CelestrakURLs struct {
	AllSatellites string `yaml:"all_satellites"`
	GeoSatellites string `yaml:"geo_satellites"`
}

//GetConfiguration From the command line arguments, get the configuration file location and parse it
func GetConfiguration() (*Info, error) {
	arg.MustParse(&args)
	fmt.Printf("Configuration file path=%s\n", args.ConfigFilePath)
	config, err := parseConfigurationFile(args.ConfigFilePath)
	return config, err
}

func parseConfigurationFile(filePath string) (*Info, error) {
	rawContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("Error when reading the configuration file: %w", err)
	}

	var c Info

	err = c.parseYaml(rawContent)
	if err != nil {
		return nil, fmt.Errorf("Error when unmarshalling the YAML configuration file: %w", err)
	}

	return &c, nil
}

func (c *Info) parseYaml(rawContent []byte) error {
	err := yaml.Unmarshal(rawContent, &c)
	if err != nil {
		return err
	}

	return nil
}
