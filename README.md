# tle-provider
Pulls TLE data from multiple sources and exposes the result through a REST API.

The OpenAPI specification is available [here.](./api/openapi-3.0.yml)

## Usage

> tle-provider -c \<configuration file>

sample configuration file:

```
celestrak_urls:
  all_satellites: "https://celestrak.com/NORAD/elements/gp.php?GROUP=active&FORMAT=json"
  geo_satellites: "https://celestrak.com/NORAD/elements/gp.php?GROUP=geo&FORMAT=json"
celestrak_refresh_rate_hours: 12
server_port: 8080
data_source: "celestrak"
```

## How to generate OpenAPI 3.0 code from spec

see https://github.com/deepmap/oapi-codegen/.

Install oapi-codegen: 

> go get github.com/deepmap/oapi-codegen/cmd/oapi-codegen

Run oapi-codegen:

- Generate server code:

> %GOPATH%/bin/oapi-codegen --config ./api/server-cfg.yml ./api/openapi-3.0.yml

- Generate types/structs: 

> %GOPATH%/bin/oapi-codegen --config ./api/types-cfg.yml ./api/openapi-3.0.yml
