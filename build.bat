%GOPATH%/bin/oapi-codegen --config ./api/server-cfg.yml ./api/openapi-3.0.yml
%GOPATH%/bin/oapi-codegen --config ./api/types-cfg.yml ./api/openapi-3.0.yml
go mod tidy
go build -o ./tle-provider.exe main.go
