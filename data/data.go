package data

type data struct {
	Results []SatelliteData `json:"results"`
}

// SatelliteData data stored in SkyminerDocument
type SatelliteData struct {
	ID         string      `json:"_id"`
	Date       string      `json:"date"`
	Satellites []Satellite `json:"satellite_data"`
}

// Satellite data structure for each satellite
type Satellite struct {
	SatelliteName string `json:"satellite_name"`
	NORADID       string `json:"norad_id"`
	TLELine1      string `json:"tle_line_1"`
	TLELine2      string `json:"tle_line_2"`
}

// Source interface for either Celestrak or Skyminer data
type Source interface {
	GetData() ([]Satellite, error)
}
