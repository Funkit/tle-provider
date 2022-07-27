package data

import (
	"github.com/go-chi/render"
	"net/http"
	"time"
)

//Info Receiving structure when parsing the configuration file
type Info struct {
	ServerPort              int                     `yaml:"server_port"`
	DataSource              string                  `yaml:"data_source"`
	CelestrakConfiguration  CelestrakConfiguration  `yaml:"celestrak_configuration"`
	FileSourceConfiguration FileSourceConfiguration `yaml:"file_source_configuration"`
}

type FileSourceConfiguration struct {
	SourceFilePath     string `yaml:"source_file_path"`
	RefreshRateSeconds int    `yaml:"refresh_rate_seconds"`
}

type CelestrakConfiguration struct {
	AllSatellitesURL string `yaml:"all_satellites_url"`
	GeoSatellitesURL string `yaml:"geo_satellites_url"`
	RefreshRateHours int    `yaml:"celestrak_refresh_rate_hours"`
}

func (i Info) IsValid() bool {
	return i.ServerPort != 0 &&
		i.DataSource != "" &&
		i.CelestrakConfiguration.RefreshRateHours != 0 &&
		i.CelestrakConfiguration.GeoSatellitesURL != "" &&
		i.CelestrakConfiguration.AllSatellitesURL != ""
}

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
