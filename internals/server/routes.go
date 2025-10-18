package server

import (
	"net/http"

	"github.com/cateruu/moto-backend/internals/middleware"
)

func (s *Server) Routes() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("GET /healthcheck", s.healthCheckHandler)

	middleware := middleware.New()
	var handler http.Handler = router
	handler = middleware.PathLogger(handler)

	return handler
}
