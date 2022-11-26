package cmd

import (
	"fmt"
	"github.com/Funkit/go-utils/utils"
	"github.com/Funkit/tle-provider/api"
	"github.com/Funkit/tle-provider/data"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/spf13/cobra"
	"log"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts a REST API server providing TLEs",
	Long: `tle-provider can get satellite TLE information
from various sources, and expose it through a single REST API.

The following sources are currently supported:
- Celestrak
- File source (following the Celestrak file formatting)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := utils.GenericYAMLParsing[data.Info](cfgFile)
		if err != nil {
			return err
		}

		if !config.IsValid() {
			panic(fmt.Errorf("invalid configuration file format"))
		}
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
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
