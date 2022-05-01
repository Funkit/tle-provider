package data

import (
	"github.com/go-chi/render"
	"net/http"
)

// Satellite data structure for each satellite
type Satellite struct {
	SatelliteName string `json:"satellite_name"`
	NORADID       int    `json:"norad_id"`
	TLELine1      string `json:"tle_line_1"`
	TLELine2      string `json:"tle_line_2"`
}

func (s *Satellite) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func GenerateRenderList(satList []*Satellite) []render.Renderer {
	var renderList []render.Renderer
	for _, sat := range satList {
		renderList = append(renderList, sat)
	}
	return renderList
}

// Source interface for either Celestrak or Skyminer data
type Source interface {
	GetData() ([]*Satellite, error)
	GetSatellite(satelliteName string) (*Satellite, error)
	GetDataSource() string
	GetConfig() (map[string]interface{}, error)
}
