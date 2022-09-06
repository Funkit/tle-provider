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
