// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.8.2 DO NOT EDIT.
package api

// Error defines model for Error.
type Error struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

// Satellite defines model for Satellite.
type Satellite struct {
	// Embedded fields due to inline allOf schema
	// Satellite name.
	Name string `json:"name"`

	// NORAD catalog ID.
	NoradId int32 `json:"norad_id"`

	// TLE line 1.
	TleLine1 string `json:"tle_line_1"`

	// TLE line 2.
	TleLine2 string `json:"tle_line_2"`
}

// ServerConfig defines model for ServerConfig.
type ServerConfig struct {
	AdditionalProperties *string `json:"additionalProperties,omitempty"`
	DataSource           string  `json:"data_source"`
}

// GetTLEListParams defines parameters for GetTLEList.
type GetTLEListParams struct {
	Group *GetTLEListParamsGroup `json:"group,omitempty"`

	// maximum number of results to return
	Limit *int32 `json:"limit,omitempty"`
}

// GetTLEListParamsGroup defines parameters for GetTLEList.
type GetTLEListParamsGroup string
