package server

import (
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	s.registerScimEndpoints(mux)
	mux.Handle("GET /", http.HandlerFunc(s.indexHandler))
	// users.AddEndpoints(mux, s.db, s.AuthenticatedMiddleware)

	// return s.corsMiddleware(s.loggingMiddleware(mux))
	return mux
}

var SCIM_PREFIX = "/scim/"

func (s *Server) registerScimEndpoints(mux *http.ServeMux) {
	
	mux.Handle("GET" + SCIM_PREFIX + "v2/Users", http.HandlerFunc(s.handleUsers))
}

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	// slog.Debug("Hello world called")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome to the SCIMtiplexer API!"))
}

func (s *Server) handleUsers(w http.ResponseWriter, r *http.Request) {
	panic("handleUsers not implemented yet")
}
