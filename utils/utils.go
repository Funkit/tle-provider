package utils

//Info Receiving structure when parsing the configuration file
type Info struct {
	ServerPort             int                    `yaml:"server_port"`
	DataSource             string                 `yaml:"data_source"`
	CelestrakConfiguration CelestrakConfiguration `yaml:"celestrak_configuration"`
}

type CelestrakConfiguration struct {
	AllSatellitesURL string `yaml:"all_satellites_url"`
	GeoSatellitesURL string `yaml:"geo_satellites_url"`
	RefreshRateHours int    `yaml:"celestrak_refresh_rate_hours"`
}

// IsALetter Check if char is an ASCII letter or not
func IsALetter(element byte) bool {
	val := ((element >= 'a') && (element <= 'z')) || ((element >= 'A') && (element <= 'Z'))

	return val
}
