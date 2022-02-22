package data

// Satellite data structure for each satellite
type Satellite struct {
	SatelliteName string `json:"satellite_name"`
	NORADID       int    `json:"norad_id"`
	TLELine1      string `json:"tle_line_1"`
	TLELine2      string `json:"tle_line_2"`
}

// Source interface for either Celestrak or Skyminer data
type Source interface {
	GetData() ([]Satellite, error)
	GetDataSource() string
	GetConfig() (map[string]string, error)
}
