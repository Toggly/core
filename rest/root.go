package rest

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Toggly/core/api"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/rs/zerolog"
)

// Server implements rest server
type Server struct {
	Version  string
	API      api.TogglyAPI
	BasePath string
	Logger   zerolog.Logger
	LogLevel zerolog.Level
}

// Run rest api
func (s *Server) Run(ctx context.Context, port int, basePath string) {
	routes := s.Router(basePath)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: chi.ServerBaseContext(ctx, routes),
	}
	go func() {
		<-ctx.Done()
		if err := srv.Shutdown(ctx); err != nil {
			s.Logger.Error().Err(err).Msg("REST stop error")
		}
		s.Logger.Info().Msg("REST server stopped")
	}()
	s.Logger.Info().Str("addr", srv.Addr).Msg("HTTP server listening")
	err := srv.ListenAndServe()
	s.Logger.Info().Msgf("HTTP server terminated, %s", err)
}

// Router returns router
func (s *Server) Router(basePath string) chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Throttle(1000))
	router.Use(middleware.Timeout(60 * time.Second))
	router.Use(middleware.Heartbeat("/ping"))
	router.Use(ServiceInfo("Toggly", s.Version))
	router.Route(basePath, s.versions)
	return router
}

func (s *Server) versions(router chi.Router) {
	router.Route("/v1", s.v1)
}

func (s *Server) v1(router chi.Router) {
	router.Use(RequestIDCtx(s.Logger))
	router.Use(Logger(s.Logger, s.LogLevel))
	router.Use(OwnerCtx(s.Logger))
	router.Use(VersionCtx("v1"))
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		log := WithRequest(s.Logger, r)
		log.Info().Msg("Some log text")
		render.PlainText(w, r, "hello")
	})
	router.Get("/nf", func(w http.ResponseWriter, r *http.Request) {
		// log := WithRequest(s.Logger, r)
		NotFoundResponse(w, r, "Did not found that")
	})
	// router.Mount("/project", (&ProjectRestAPI{API: s.API}).Routes())
	// router.Mount("/project/{project_code}/env", (&EnvironmentRestAPI{API: s.API}).Routes())
	// router.Mount("/project/{project_code}/env/{env_code}/object", (&ObjectRestAPI{API: s.API}).Routes())
}

func owner(s *http.Request) string {
	return OwnerFromContext(s)
}

func projectCode(s *http.Request) string {
	return chi.URLParam(s, "project_code")
}

func environmentCode(s *http.Request) string {
	return chi.URLParam(s, "env_code")
}

func objectCode(s *http.Request) string {
	return chi.URLParam(s, "object_code")
}
