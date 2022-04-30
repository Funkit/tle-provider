package main

import (
	"github.com/Funkit/tle-provider/api"
	"log"

	"github.com/Funkit/tle-provider/data"
	"github.com/Funkit/tle-provider/utils"
)

var args struct {
	ConfigFilePath string `arg:"required,-c" help:"path to the configuration file"`
}

var Version = "development"

func main() {
	log.Println("Version:\t", Version)

	config, err := utils.GenericYAMLParsing[utils.Info](args.ConfigFilePath)
	if err != nil {
		panic(err)
	}

	var source data.Source

	switch config.DataSource {
	case "celestrak":
		source = data.NewCelestrakClient(config.CelestrakConfiguration.AllSatellitesURL, config.CelestrakConfiguration.GeoSatellitesURL, config.CelestrakConfiguration.RefreshRateHours)
	}

	server := api.NewServer(config.ServerPort, source)

	if err := server.Run(); err != nil {
		panic(err)
	}
}
