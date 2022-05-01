package main

import (
	"github.com/Funkit/tle-provider/api"
	"github.com/Funkit/tle-provider/data"
	"github.com/alexflint/go-arg"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log"

	"github.com/Funkit/tle-provider/utils"
)

var args struct {
	ConfigFilePath string `arg:"required,-c" help:"path to the configuration file"`
}

var Version = "development"

func main() {
	log.Println("Version:\t", Version)

	arg.MustParse(&args)

	config, err := utils.GenericYAMLParsing[utils.Info](args.ConfigFilePath)
	if err != nil {
		panic(err)
	}

	var source data.Source

	switch config.DataSource {
	case "celestrak":
		source = data.NewCelestrakClient(
			config.CelestrakConfiguration.AllSatellitesURL,
			config.CelestrakConfiguration.GeoSatellitesURL,
			config.CelestrakConfiguration.RefreshRateHours)
	}

	server := api.NewServer(config.ServerPort, source)
	server.AddMiddlewares(middleware.Logger, render.SetContentType(render.ContentTypeJSON), middleware.Recoverer)
	server.InitializeRoutes()

	if err := server.Run(); err != nil {
		panic(err)
	}
}
