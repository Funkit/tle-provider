package data

import (
	"github.com/go-chi/render"
	"net/http"
	"time"
)

// Satellite data structure for each satellite
type Satellite struct {
	SatelliteName string `json:"satellite_name"`
	NORADID       int    `json:"norad_id"`
	TLELine1      string `json:"tle_line_1"`
	TLELine2      string `json:"tle_line_2"`
}

type SatelliteErr struct {
	Err error
	Sat Satellite
}

func (s Satellite) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s Satellite) IsNull() bool {
	return len(s.TLELine1) != 69
}

func GenerateRenderList(satList []Satellite) []render.Renderer {
	var renderList []render.Renderer
	for _, sat := range satList {
		renderList = append(renderList, sat)
	}
	return renderList
}

// Source interface for either Celestrak or Skyminer data
type Source interface {
	Update(done <-chan struct{}, period time.Duration)
	GetData() ([]Satellite, error)
	GetSatellite(satelliteName string) chan SatelliteErr
	GetDataSource() string
	GetConfig() (map[string]interface{}, error)
}
