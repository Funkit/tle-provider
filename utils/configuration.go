package utils

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/alexflint/go-arg"
	"gopkg.in/yaml.v2"
)

var args struct {
	ConfigFilePath string `arg:"required,-c" help:"path to the configuration file"`
}

// GetConfiguration From the command line arguments, get the configuration file location and parse it
func GetConfiguration() (map[string]interface{}, error) {
	arg.MustParse(&args)
	fmt.Printf("Configuration file path=%s\n", args.ConfigFilePath)
	config, err := parseConfigurationFile(args.ConfigFilePath)
	return config, err
}

func parseConfigurationFile(filePath string) (map[string]interface{}, error) {
	rawContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("Error when reading the configuration file: %w", err)
	}

	var config map[string]interface{}

	err = yaml.Unmarshal(rawContent, &config)
	if err != nil {
		return nil, fmt.Errorf("Error when unmarshalling the YAML configuration file: %w", err)
	}

	// Check that at least server port and data source are defined in the configuration file
	_, portIsInt := config["server_port"].(int)
	if !portIsInt {
		return nil, errors.New("Server port is of the wrong format (should be int)")
	}
	_, dataSourceIsString := config["data_source"].(string)
	if !dataSourceIsString {
		return nil, errors.New("Data source is of the wrong format (should be string)")
	}

	return config, nil
}
