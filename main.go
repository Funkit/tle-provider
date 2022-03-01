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
func NewDataSource(info map[string]interface{}) (data.Source, error) {

	if info["data_source"] == "celestrak" {
		return data.NewCelestrakClient(info)
	} else if info["data_source"] == "postgresql" {
		return data.NewPostgreSQLClient(info)
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

	tleServer := api.NewTLEServer(dataSource)

	// Echo router setup
	e := echo.New()
	e.Use(echomiddleware.Logger())
	e.Use(middleware.OapiRequestValidator(swagger))
	api.RegisterHandlers(e, tleServer)

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf("0.0.0.0:%d", config["server_port"])))
}
