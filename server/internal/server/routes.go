package server

import (
	"log/slog"
	"net/http"
	scimuser "github.com/jawee/scimtiplexer/internal/scim/user"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	// s.registerScimEndpoints(mux)

	scimuser.RegisterEndpoints(mux, s.db.GetRepository())

	return s.corsMiddleware(s.loggingMiddleware(mux))
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Replace "*" with specific origins if needed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token, ApiKey")
		w.Header().Set("Access-Control-Allow-Credentials", "false") // Set to "true" if credentials are required

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			slog.Debug("Returning 204")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		slog.Debug("Proceeding")
		// Proceed with the next handler
		next.ServeHTTP(w, r)
	})
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("Request", "Path", r.URL.Path, "Method", r.Method)

		// Proceed with the next handler
		next.ServeHTTP(w, r)
	})
}


