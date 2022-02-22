package api

import (
	"encoding/json"
	"net/http"

	"github.com/Funkit/tle-provider/data"
	"github.com/labstack/echo/v4"
)

// TLEServer The TLE server object
type TLEServer struct {
	Source data.Source
}

// NewTLEServer Generates a new server
func NewTLEServer(s data.Source) *TLEServer {
	return &TLEServer{
		Source: s,
	}
}

func sendServerError(ctx echo.Context, code int, message string) error {
	petErr := Error{
		Code:    int32(code),
		Message: message,
	}
	err := ctx.JSON(code, petErr)
	return err
}

// GetTLEList returns all the latest available TLEs
func (ts *TLEServer) GetTLEList(ctx echo.Context, params GetTLEListParams) error {

	satelliteList, err := ts.Source.GetData()
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, satelliteList)
}

// FindASatelliteByName Return TLE for a given satellite
func (ts *TLEServer) FindASatelliteByName(ctx echo.Context, satellite string) error {

	satelliteList, err := ts.Source.GetData()
	if err != nil {
		return err
	}

	for _, sat := range satelliteList {
		if sat.SatelliteName == satellite {
			return ctx.JSON(http.StatusOK, sat)
		}
	}

	return sendServerError(ctx, http.StatusBadRequest, "Satellite not found")
}

// GetConfig Return TLE for a given satellite
func (ts *TLEServer) GetConfig(ctx echo.Context) error {

	conf, err := ts.Source.GetConfig()
	if err != nil {
		return err
	}
	confJSON, err := json.Marshal(conf)
	if err != nil {
		return err
	}
	confJSONString := string(confJSON)

	return ctx.JSON(http.StatusOK, ServerConfig{
		DataSource:           ts.Source.GetDataSource(),
		AdditionalProperties: &confJSONString,
	})
}
