[![Go Report Card](https://goreportcard.com/badge/github.com/Funkit/tle-provider)](https://goreportcard.com/report/github.com/Funkit/tle-provider)

# TLE Provider
Pulls TLE data from multiple sources and exposes the result through a REST API.

The OpenAPI specification is available [here.](./api/openapi-3.0.yml)

The following data sources are available :

- **Celestrak:** Query data directly from the Celestrak JSON API URLs
- **File:** Expose the data pulled from the text file dump generated from Celestrak (example available in the `samples` folder).

## Usage

> tle-provider serve --config \<configuration file>

sample configuration file:

```
server_port: 5000
data_source: "celestrak"
celestrak_configuration:
  all_satellites_url: "https://celestrak.com/NORAD/elements/gp.php?GROUP=active&FORMAT=json"
  geo_satellites_url: "https://celestrak.com/NORAD/elements/gp.php?GROUP=geo&FORMAT=json"
  celestrak_refresh_rate_hours: 12
file_source_configuration:
  source_file_path: "./samples/active_satellites_tle.txt"
  refresh_rate_seconds: 30
```

- `server_port`: exposed port for the service.
- `data_source`: either `celestrak` or `file`.
- `celestrak_configuration`:
  - `all_satellites_url`: the URL to use when querying the Celestrak website for all satellites.
  - `geo_satellites_url` the URL to use when querying the Celestrak website for geosynchrnous satellites only.
  - `celestrak_refresh_rate_hours`: period at which to query the data from Celestrak.
- `file_source_configuration`:
  - `source_file_path`: path to the TLE source file.
  - `refresh_rate_seconds`: revisit rate of the source file.

**Note**: when performing `Run()`, the server starts a separate thread for pulling data from the source only if the refresh rate is set at more than 1 second.
