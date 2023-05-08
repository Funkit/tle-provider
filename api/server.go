package api

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sync"
	"time"

	"github.com/Funkit/go-utils/apierror"
	"github.com/Funkit/tle-provider/data"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

var (
	Constellations = map[string]*regexp.Regexp{
		"oneweb":   regexp.MustCompile("ONEWEB-[0-9]+"),
		"starlink": regexp.MustCompile("STARLINK-[0-9]+"),
	}
)

type Server struct {
	source                 data.Source
	router                 chi.Router
	Port                   int
	DataRefreshRate        time.Duration
	CelestrakRefreshRate   time.Duration
	FileRefreshRateSeconds time.Duration
	mu                     sync.RWMutex
	satellitesTLEs         []data.Satellite
	satellitesTLEsMap      map[string]data.Satellite
	constellationsTLEs     map[string][]data.Satellite
	lastPull               time.Time
	done                   chan struct{}
}

func NewServer(port int, source data.Source, refreshRate time.Duration) *Server {
	done := make(chan struct{})
	return &Server{
		source:          source,
		router:          chi.NewRouter(),
		Port:            port,
		DataRefreshRate: refreshRate,
		done:            done,
		lastPull:        time.Date(1970, 01, 01, 0, 0, 0, 1, time.UTC),
	}
}

func (s *Server) AddMiddlewares(middlewares ...func(handler http.Handler) http.Handler) {
	s.router.Use(middlewares...)
}

func (s *Server) SubRoutes(baseURL string, r chi.Router) {
	s.router.Mount(baseURL, r)
}

func (s *Server) update() {
	for {
		select {
		case <-s.done:
			log.Println("END")
			break
		case <-time.After(s.DataRefreshRate):
			sats, err := s.source.GetData()
			if err != nil {
				log.Println(err.Error())
			} else {
				s.UpdateAllValues(sats)
			}
		}
	}
}

func (s *Server) UpdateAllValues(sats []data.Satellite) {
	s.mu.Lock()
	s.satellitesTLEs = sats
	s.satellitesTLEsMap = make(map[string]data.Satellite)
	s.constellationsTLEs = make(map[string][]data.Satellite)

	for _, element := range sats {
		s.satellitesTLEsMap[element.SatelliteName] = element

		for constName, namePattern := range Constellations {
			if namePattern.MatchString(element.SatelliteName) {
				s.constellationsTLEs[constName] = append(s.constellationsTLEs[constName], element)
			}
		}
	}

	s.lastPull = time.Now()
	log.Printf("data successfully pulled from %s at %s\n", s.source.GetDataSource(), time.Now().Format("2006-01-02T15:04:05Z"))
	s.mu.Unlock()
}

func (s *Server) Run() error {

	sats, err := s.source.GetData()
	if err != nil {
		return err
	} else {
		s.UpdateAllValues(sats)
	}

	if s.DataRefreshRate >= time.Second {
		go s.update()
	}

	log.Printf("Listening on port %v\n", s.Port)

	if err := http.ListenAndServe(fmt.Sprintf(":%v", s.Port), s.router); err != nil {
		s.done <- struct{}{}
		panic(err)
	}

	s.done <- struct{}{}
	return nil
}

func (s *Server) InitializeRoutes() {
	s.router.Get("/tle", s.getTLEList())
	s.router.Get("/tle/{satellite}", s.getTLE())
}

func (s *Server) getTLEList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		constellation, ok := r.URL.Query()["constellation"]
		if ok && len(constellation) != 0 {
			if len(constellation[0]) != 0 {

				s.mu.RLock()
				defer s.mu.RUnlock()
				if len(s.constellationsTLEs[constellation[0]]) == 0 {
					apierror.Handle(w, r, apierror.Wrap(fmt.Errorf("no satellites found"), apierror.ErrNotFound))
					return
				}

				renderList := data.GenerateRenderList(s.constellationsTLEs[constellation[0]])
				if err := render.RenderList(w, r, renderList); err != nil {
					apierror.Handle(w, r, apierror.Wrap(err, apierror.ErrRender))
				}
				return
			}
		}

		s.mu.RLock()
		defer s.mu.RUnlock()
		if len(s.satellitesTLEs) == 0 {
			apierror.Handle(w, r, apierror.Wrap(fmt.Errorf("no satellite found"), apierror.ErrNotFound))
			return
		}

		renderList := data.GenerateRenderList(s.satellitesTLEs)
		if err := render.RenderList(w, r, renderList); err != nil {
			apierror.Handle(w, r, apierror.Wrap(err, apierror.ErrRender))
		}
	}
}

func (s *Server) getTLE() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		satelliteParam := chi.URLParam(r, "satellite")

		s.mu.RLock()
		defer s.mu.RUnlock()

		if s.satellitesTLEsMap[satelliteParam].IsNull() {
			apierror.Handle(w, r, apierror.Wrap(fmt.Errorf("satellite %v not found", satelliteParam), apierror.ErrNotFound))
			return
		}

		if err := render.Render(w, r, s.satellitesTLEsMap[satelliteParam]); err != nil {
			apierror.Handle(w, r, apierror.Wrap(err, apierror.ErrRender))
		}
	}
}
