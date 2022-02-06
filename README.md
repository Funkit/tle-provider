# tle-provider
Pulls TLE data from multiple sources and exposes the result through a REST API

## How to generate OpenAPI 3.0 code from spec

see https://github.com/deepmap/oapi-codegen/.

Install oapi-codegen: 

> go get github.com/deepmap/oapi-codegen/cmd/oapi-codegen

Run oapi-codegen:

- Generate server code:

> %GOPATH%/bin/oapi-codegen --config ./api/server-cfg.yml ./api/openapi-3.0.yml

- Generate types/structs: 

> %GOPATH%/bin/oapi-codegen --config ./api/types-cfg.yml ./api/openapi-3.0.yml
