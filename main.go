package main

import (
	"errors"
	"fmt"

	"github.com/Funkit/tle-provider/api"
	"github.com/Funkit/tle-provider/data"
	"github.com/Funkit/tle-provider/utils"
	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

// NewDataSource Create data source
func NewDataSource(info *utils.Info) (data.Source, error) {

	if info.DataSource == "celestrak" {
		s := data.NewCelestrakClient(info)
		return s, nil
	}
	return nil, errors.New("data source not supported")
}

func main() {
	config, err := utils.GetConfiguration()
	if err != nil {
		panic(err)
	}

	// TLE server setup
	swagger, err := api.GetSwagger()
	if err != nil {
		panic(fmt.Errorf("Error loading swagger spec\n: %s", err))
	}
	swagger.Servers = nil

	dataSource, err := NewDataSource(config)
	if err != nil {
		panic(err)
	}

	tleServer := api.NewTLEServer(config, dataSource)

	// Echo router setup
	e := echo.New()
	e.Use(echomiddleware.Logger())
	e.Use(middleware.OapiRequestValidator(swagger))
	api.RegisterHandlers(e, tleServer)

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf("0.0.0.0:%d", config.ServerPort)))
}
