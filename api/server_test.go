package api

import (
	"fmt"
	"github.com/Funkit/tle-provider/data"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func executeRequest(req *http.Request, s *Server) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)

	return rr
}

func TestGetTLE(t *testing.T) {
	done := make(chan struct{})

	source := data.NewFileSource(
		"../samples/tle.txt",
		30)

	source.Update(done, 30*time.Second)
	defer func() { done <- struct{}{} }()

	type fields struct {
		satelliteName string
	}
	tests := []struct {
		name         string
		fields       fields
		wantRespCode int
		wantBody     string
	}{
		{
			name: "Working case",
			fields: fields{
				satelliteName: "LAGEOS%201",
			},
			wantRespCode: http.StatusOK,
			wantBody:     "{\"satellite_name\":\"LAGEOS 1\",\"norad_id\":8820,\"tle_line_1\":\"1 08820U 76039A   22206.68532073  .00000028  00000-0  00000-0 0  9999\",\"tle_line_2\":\"2 08820 109.8533  52.0899 0045094 246.5947 308.4924  6.38664901822297\"}\n",
		},
		{
			name: "Satellite not found",
			fields: fields{
				satelliteName: "THISSATELLITEDOESNOTEXIST",
			},
			wantRespCode: http.StatusNotFound,
			wantBody:     "{\"status\":404,\"message\":\"Resource Not found\"}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewServer(80, source, 1, 30)
			s.AddMiddlewares(middleware.Logger, render.SetContentType(render.ContentTypeJSON), middleware.Recoverer)
			s.InitializeRoutes()

			req, _ := http.NewRequest("GET", fmt.Sprintf("/tle/%v", tt.fields.satelliteName), nil)

			response := executeRequest(req, s)
			if response.Code != tt.wantRespCode {
				t.Errorf("Expected response code %d. Got %d\n", tt.wantRespCode, response.Code)
			} else {
				t.Log(len(tt.wantBody))
				t.Log(len(response.Body.String()))
				if response.Body.String() != tt.wantBody {
					t.Errorf("Expected response body %s. Got %s\n", tt.wantBody, response.Body.String())
				}
			}
		})
	}
}

func TestGetTLEList(t *testing.T) {
	done := make(chan struct{})

	source := data.NewFileSource(
		"../samples/tle.txt",
		30)

	source.Update(done, 30*time.Second)
	defer func() { done <- struct{}{} }()

	type fields struct {
		constellation string
	}
	tests := []struct {
		name         string
		fields       fields
		wantRespCode int
		wantBody     string
	}{
		{
			name:         "Get all TLEs",
			fields:       fields{},
			wantRespCode: http.StatusOK,
			wantBody: "[" +
				"{\"satellite_name\":\"OPS 5712 (P/L 153)\",\"norad_id\":2874,\"tle_line_1\":\"1 02874U 67053H   22206.60472723 -.00000017  00000-0  26447-4 0  9991\",\"tle_line_2\":\"2 02874  69.9738 283.4261 0009834 250.7192 109.2850 13.96410943808158\"}," +
				"{\"satellite_name\":\"CALSPHERE 1\",\"norad_id\":900,\"tle_line_1\":\"1 00900U 64063C   22206.83199285  .00000371  00000-0  38562-3 0  9993\",\"tle_line_2\":\"2 00900  90.1732  41.6116 0024844 266.8448 104.5887 13.73849434875933\"}," +
				"{\"satellite_name\":\"LAGEOS 1\",\"norad_id\":8820,\"tle_line_1\":\"1 08820U 76039A   22206.68532073  .00000028  00000-0  00000-0 0  9999\",\"tle_line_2\":\"2 08820 109.8533  52.0899 0045094 246.5947 308.4924  6.38664901822297\"}," +
				"{\"satellite_name\":\"ONEWEB-0012\",\"norad_id\":44057,\"tle_line_1\":\"1 44057U 19010A   22206.81764082 -.00000043  00000+0 -14585-3 0  9993\",\"tle_line_2\":\"2 44057  87.9150 151.8950 0002369 106.7932 253.3459 13.16592117164401\"}," +
				"{\"satellite_name\":\"ONEWEB-0010\",\"norad_id\":44058,\"tle_line_1\":\"1 44058U 19010B   22206.61495996  .00000036  00000+0  59419-4 0  9993\",\"tle_line_2\":\"2 44058  87.9152 151.9295 0002498  91.0942 269.0475 13.16593199164422\"}," +
				"{\"satellite_name\":\"ONEWEB-0008\",\"norad_id\":44059,\"tle_line_1\":\"1 44059U 19010C   22206.56425428 -.00000049  00000+0 -16274-3 0  9998\",\"tle_line_2\":\"2 44059  87.9155 151.9775 0001572  76.2188 283.9118 13.16592146164530\"}," +
				"{\"satellite_name\":\"STARLINK-61\",\"norad_id\":44249,\"tle_line_1\":\"1 44249U 19029Q   22207.21981331  .01879246  22465-2  37398-2 0  9995\",\"tle_line_2\":\"2 44249  52.9518 229.8866 0008132  34.7714 325.3242 15.99740001176748\"}," +
				"{\"satellite_name\":\"STARLINK-71\",\"norad_id\":44252,\"tle_line_1\":\"1 44252U 19029T   22206.63642171  .00063016  00000+0  13933-2 0  9996\",\"tle_line_2\":\"2 44252  52.9947 285.2994 0003334  27.1844 332.9334 15.43254345174817\"}" +
				"]\n",
		},
		{
			name: "Get OneWeb constellation",
			fields: fields{
				constellation: "oneweb",
			},
			wantRespCode: http.StatusOK,
			wantBody: "[" +
				"{\"satellite_name\":\"ONEWEB-0012\",\"norad_id\":44057,\"tle_line_1\":\"1 44057U 19010A   22206.81764082 -.00000043  00000+0 -14585-3 0  9993\",\"tle_line_2\":\"2 44057  87.9150 151.8950 0002369 106.7932 253.3459 13.16592117164401\"}," +
				"{\"satellite_name\":\"ONEWEB-0010\",\"norad_id\":44058,\"tle_line_1\":\"1 44058U 19010B   22206.61495996  .00000036  00000+0  59419-4 0  9993\",\"tle_line_2\":\"2 44058  87.9152 151.9295 0002498  91.0942 269.0475 13.16593199164422\"}," +
				"{\"satellite_name\":\"ONEWEB-0008\",\"norad_id\":44059,\"tle_line_1\":\"1 44059U 19010C   22206.56425428 -.00000049  00000+0 -16274-3 0  9998\",\"tle_line_2\":\"2 44059  87.9155 151.9775 0001572  76.2188 283.9118 13.16592146164530\"}" +
				"]\n",
		},
		{
			name: "Get Starlink constellation",
			fields: fields{
				constellation: "starlink",
			},
			wantRespCode: http.StatusOK,
			wantBody: "[" +
				"{\"satellite_name\":\"STARLINK-61\",\"norad_id\":44249,\"tle_line_1\":\"1 44249U 19029Q   22207.21981331  .01879246  22465-2  37398-2 0  9995\",\"tle_line_2\":\"2 44249  52.9518 229.8866 0008132  34.7714 325.3242 15.99740001176748\"}," +
				"{\"satellite_name\":\"STARLINK-71\",\"norad_id\":44252,\"tle_line_1\":\"1 44252U 19029T   22206.63642171  .00063016  00000+0  13933-2 0  9996\",\"tle_line_2\":\"2 44252  52.9947 285.2994 0003334  27.1844 332.9334 15.43254345174817\"}" +
				"]\n",
		},
		{
			name: "wrong constellation argument",
			fields: fields{
				constellation: "DOESNOTEXIST",
			},
			wantRespCode: http.StatusNotFound,
			wantBody:     "{\"status\":404,\"message\":\"Resource Not found\"}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewServer(80, source, 1, 30)
			s.AddMiddlewares(middleware.Logger, render.SetContentType(render.ContentTypeJSON), middleware.Recoverer)
			s.InitializeRoutes()

			path := "/tle"

			if tt.fields.constellation != "" {
				path = fmt.Sprintf("%s?constellation=%s", path, tt.fields.constellation)
			}

			req, _ := http.NewRequest("GET", path, nil)

			response := executeRequest(req, s)
			if response.Code != tt.wantRespCode {
				t.Errorf("Expected response code %d. Got %d\n", tt.wantRespCode, response.Code)
			} else {
				t.Log(response.Body.String())
				t.Log(len(tt.wantBody))
				if response.Body.String() != tt.wantBody {
					t.Errorf("Expected response body %s. Got %s\n", tt.wantBody, response.Body.String())
				}
			}
		})
	}
}
