%GOPATH%/bin/oapi-codegen --config ./api/server-cfg.yml ./api/openapi-3.0.yml
%GOPATH%/bin/oapi-codegen --config ./api/types-cfg.yml ./api/openapi-3.0.yml
go build -o ./tleprovider.exe main.go
