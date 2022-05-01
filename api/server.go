package api

import (
	"fmt"
	"github.com/Funkit/tle-provider/apierror"
	"github.com/Funkit/tle-provider/data"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log"
	"net/http"
	"sync"
)

type Server struct {
	source data.Source
	router chi.Router
	Port   int
	mu     sync.RWMutex
}

func NewServer(port int, source data.Source) *Server {
	return &Server{
		source: source,
		router: chi.NewRouter(),
		Port:   port,
	}
}

func (s *Server) AddMiddlewares(middlewares ...func(handler http.Handler) http.Handler) {
	s.router.Use(middlewares...)
}

func (s *Server) SubRoutes(baseURL string, r chi.Router) {
	s.router.Mount(baseURL, r)
}

func (s *Server) Run() error {
	log.Printf("Listening on port %v\n", s.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%v", s.Port), s.router); err != nil {
		panic(err)
	}
	return nil
}

func (s *Server) InitializeRoutes() {
	s.router.Get("/tle", s.getTLEList())
	s.router.Get("/tle/{satellite}", s.getTLE())
}

func (s *Server) getTLEList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		satelliteList, err := s.source.GetData()
		if err != nil {
			apierror.Handle(w, r, err)
			return
		}

		renderList := data.GenerateRenderList(satelliteList)
		if err := render.RenderList(w, r, renderList); err != nil {
			apierror.Handle(w, r, apierror.Wrap(err, apierror.ErrRender))
		}
	}
}

func (s *Server) getTLE() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		satelliteParam := chi.URLParam(r, "satellite")
		satellite, err := s.source.GetSatellite(satelliteParam)
		if err != nil {
			apierror.Handle(w, r, err)
			return
		}

		if err := render.Render(w, r, satellite); err != nil {
			apierror.Handle(w, r, apierror.Wrap(err, apierror.ErrRender))
		}
	}
}
