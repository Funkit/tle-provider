package main

import (
	"fmt"
	"github.com/Funkit/go-utils/utils"
	"github.com/Funkit/tle-provider/api"
	"github.com/Funkit/tle-provider/data"
	"github.com/alexflint/go-arg"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log"
)

var args struct {
	ConfigFilePath string `arg:"required,-c" help:"path to the configuration file"`
	Demo           bool   `default:"false"`
}

var Version = "development"

func main() {

	arg.MustParse(&args)

	config, err := utils.GenericYAMLParsing[data.Info](args.ConfigFilePath)
	if err != nil {
		panic(err)
	}

	if !config.IsValid() {
		panic(fmt.Errorf("invalid configuration file format"))
	}

	log.Println("Version:\t ", Version)
	log.Println("Data source: ", config.DataSource)

	var source data.Source

	switch config.DataSource {
	case "celestrak":
		source = data.NewCelestrakClient(
			config.CelestrakConfiguration.AllSatellitesURL,
			config.CelestrakConfiguration.GeoSatellitesURL,
			config.CelestrakConfiguration.RefreshRateHours)
	case "file":
		source = data.NewFileSource(
			config.FileSourceConfiguration.SourceFilePath,
			config.FileSourceConfiguration.RefreshRateSeconds)
	}

	server := api.NewServer(config.ServerPort, source, config.CelestrakConfiguration.RefreshRateHours, config.FileSourceConfiguration.RefreshRateSeconds)
	server.AddMiddlewares(middleware.Logger, render.SetContentType(render.ContentTypeJSON), middleware.Recoverer)
	server.InitializeRoutes()

	if err := server.Run(); err != nil {
		panic(err)
	}
}
