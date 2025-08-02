package server

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"strings"
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
	
	mux.Handle("GET " + SCIM_PREFIX + "v2/users", s.ScimEndpointAuth(http.HandlerFunc(s.handleUsers)))
}

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	// slog.Debug("Hello world called")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome to the SCIMtiplexer API!"))
}

func (s *Server) handleUsers(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handleUsers called for organisation", "orgid", r.Context().Value("orgid"))
}

func (s *Server) ScimEndpointAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		repo := s.db.GetRepository();

		slog.Debug("ScimEndpointAuth called", "method", r.Method, "url", r.URL.Path)

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		//TODO: Check if token matches in DB and set org id in context to be used later
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := repo.GetOrganisationTokenByToken(r.Context(), tokenStr)
		if err != nil {
			slog.Error("GetOrganisationTokenByToken failed", "error", err)
			if err == sql.ErrNoRows {
				slog.Info("Token not found in database", "token", tokenStr)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		claimsCtx := context.WithValue(r.Context(), "orgid", token.OrganisationID)
		r = r.WithContext(claimsCtx)
		
		next.ServeHTTP(w, r)
	})
}

