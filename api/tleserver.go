package api

import (
	"net/http"

	"github.com/Funkit/tle-provider/data"
	"github.com/Funkit/tle-provider/utils"
	"github.com/labstack/echo/v4"
)

//TLEServer The TLE server object
type TLEServer struct {
	Source data.Source
}

//NewTLEServer Generates a new server
func NewTLEServer(config *utils.Info, dsBuilder data.SourceBuilder) (*TLEServer, error) {

	s, err := dsBuilder.NewDataSource(config.DataSource)
	if err != nil {
		return nil, err
	}

	return &TLEServer{
		Source: s,
	}, nil
}

func sendServerError(ctx echo.Context, code int, message string) error {
	petErr := Error{
		Code:    int32(code),
		Message: message,
	}
	err := ctx.JSON(code, petErr)
	return err
}

//GetTLEList returns all the latest available TLEs
func (ts *TLEServer) GetTLEList(ctx echo.Context, params GetTLEListParams) error {

	satelliteList, err := ts.Source.GetData()
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, satelliteList)
}

//FindASatelliteByName Return TLE for a given satellite
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
