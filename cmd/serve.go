package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/Funkit/go-utils/utils"
	"github.com/Funkit/tle-provider/api"
	"github.com/Funkit/tle-provider/data"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/spf13/cobra"
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

		var refreshRate time.Duration

		switch config.DataSource {
		case "celestrak":
			source = data.NewCelestrakClient(
				config.CelestrakConfiguration.AllSatellitesURL,
				config.CelestrakConfiguration.GeoSatellitesURL)
			refreshRate = time.Duration(config.CelestrakConfiguration.RefreshRateHours) * time.Hour
		case "file":
			source = data.NewFileSource(
				config.FileSourceConfiguration.SourceFilePath)
			refreshRate = time.Duration(config.FileSourceConfiguration.RefreshRateSeconds) * time.Second
		}

		server := api.NewServer(config.ServerPort, source, refreshRate)
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
