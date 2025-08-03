package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	s.registerScimEndpoints(mux)
	// users.AddEndpoints(mux, s.db, s.AuthenticatedMiddleware)

	// return s.corsMiddleware(s.loggingMiddleware(mux))
	return s.loggingMiddleware(mux)
}

var SCIM_PREFIX = "/scim/v2/"

func (s *Server) registerScimEndpoints(mux *http.ServeMux) {
	s.registerScimEndpoint(mux, "GET", "Users", http.HandlerFunc(s.handleGetUsers))
	s.registerScimEndpoint(mux, "GET", "Users/", http.HandlerFunc(s.handleGetUsers))
	s.registerScimEndpoint(mux, "POST", "Users", http.HandlerFunc(s.handlePostUsers))

	s.registerScimEndpoint(mux, "GET", "Users/{id}", http.HandlerFunc(s.handleGetUsersById))
}

func (s *Server) registerScimEndpoint(mux *http.ServeMux, method, resource string, handler http.Handler) {
	mux.Handle(method+" "+SCIM_PREFIX+resource, s.ScimEndpointAuth(handler))
	mux.Handle(method+" "+SCIM_PREFIX+strings.ToLower(resource), s.ScimEndpointAuth(handler))
}

func (s *Server) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handleGetUsers called for organisation", "orgid", r.Context().Value("orgid"))
	queryParams := r.URL.Query()

	slog.Debug("Query parameters", "params", queryParams)

	// Here you would typically fetch users from the database based on the organisation ID
	// For now, we will just return a placeholder response
	w.Header().Set("Content-Type", "application/scim+json")
	w.WriteHeader(http.StatusOK)
	user := DummySCIMUser("")
	userResp := NewSCIMUserListResponse([]User{user}, 1, 1, 1)
	jsonOutput, _ := json.Marshal(userResp)
	w.Write(jsonOutput)
}

func (s *Server) handlePostUsers(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handlePostUsers called for organisation", "orgid", r.Context().Value("orgid"))

	var userReq UserCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
		slog.Error("Failed to decode user creation request", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	slog.Debug("User creation request", "request", userReq)

	w.Header().Set("Content-Type", "application/scim+json")
	w.WriteHeader(http.StatusCreated)
	userResp := DummySCIMUser("")
	jsonOutput, _ := json.Marshal(userResp)
	w.Write(jsonOutput)
}

func (s *Server) handleGetUsersById(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handleGetUsersById called for organisation", "orgid", r.Context().Value("orgid"))
	requestedId := r.PathValue("id")
	slog.Debug("Requested user ID", "id", requestedId)

	w.Header().Set("Content-Type", "application/scim+json")
	w.WriteHeader(http.StatusOK)
	userResp := DummySCIMUser("")
	jsonOutput, _ := json.Marshal(userResp)
	w.Write(jsonOutput)
}

func (s *Server) ScimEndpointAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		repo := s.db.GetRepository()

		slog.Debug("ScimEndpointAuth called", "method", r.Method, "url", r.URL.Path)

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

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

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("Request", "Path", r.URL.Path, "Method", r.Method)

		// Proceed with the next handler
		next.ServeHTTP(w, r)
	})
}

const (
	SchemaUser                  = "urn:ietf:params:scim:schemas:core:2.0:User"
	SchemaGroup                 = "urn:ietf:params:scim:schemas:core:2.0:Group"
	SchemaEnterpriseUser        = "urn:ietf:params:scim:schemas:extension:enterprise:2.0:User"
	SchemaServiceProviderConfig = "urn:ietf:params:scim:schemas:core:2.0:ServiceProviderConfig"
	SchemaResourceType          = "urn:ietf:params:scim:schemas:core:2.0:ResourceType"
	SchemaSchema                = "urn:ietf:params:scim:schemas:core:2.0:Schema"
	SchemaListResponse          = "urn:ietf:params:scim:api:messages:2.0:ListResponse"
	SchemaError                 = "urn:ietf:params:scim:api:messages:2.0:Error"
)

type Meta struct {
	ResourceType string    `json:"resourceType"`
	Created      time.Time `json:"created"`
	LastModified time.Time `json:"lastModified"`
	Location     string    `json:"location"`
	Version      string    `json:"version,omitempty"`
}

type Name struct {
	Formatted       string `json:"formatted,omitempty"`
	FamilyName      string `json:"familyName,omitempty"`
	GivenName       string `json:"givenName,omitempty"`
	MiddleName      string `json:"middleName,omitempty"`
	HonorificPrefix string `json:"honorificPrefix,omitempty"`
	HonorificSuffix string `json:"honorificSuffix,omitempty"`
}

type Email struct {
	Value   string `json:"value"`
	Display string `json:"display,omitempty"`
	Type    string `json:"type,omitempty"`
	Primary bool   `json:"primary,omitempty"`
}

type PhoneNumber struct {
	Value   string `json:"value"`
	Display string `json:"display,omitempty"`
	Type    string `json:"type,omitempty"`
	Primary bool   `json:"primary,omitempty"`
}

type GroupMember struct {
	Value   string `json:"value"`
	Ref     string `json:"$ref,omitempty"`
	Display string `json:"display,omitempty"`
	Type    string `json:"type,omitempty"`
}

type Manager struct {
	Value       string `json:"value"`
	Ref         string `json:"$ref,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}

type EnterpriseUserExtension struct {
	EmployeeNumber string   `json:"employeeNumber,omitempty"`
	Organization   string   `json:"organization,omitempty"`
	Department     string   `json:"department,omitempty"`
	Division       string   `json:"division,omitempty"`
	CostCenter     string   `json:"costCenter,omitempty"`
	Manager        *Manager `json:"manager,omitempty"`
}

type User struct {
	Schemas    []string `json:"schemas"`
	ID         string   `json:"id"`
	ExternalID string   `json:"externalId,omitempty"`
	Meta       Meta     `json:"meta"`

	UserName          string `json:"userName"`
	Name              *Name  `json:"name,omitempty"`
	DisplayName       string `json:"displayName,omitempty"`
	NickName          string `json:"nickName,omitempty"`
	ProfileURL        string `json:"profileUrl,omitempty"`
	Title             string `json:"title,omitempty"`
	UserType          string `json:"userType,omitempty"`
	PreferredLanguage string `json:"preferredLanguage,omitempty"`
	Locale            string `json:"locale,omitempty"`
	Timezone          string `json:"timezone,omitempty"`
	Active            bool   `json:"active"`
	Password          string `json:"password,omitempty"`

	Emails       []Email       `json:"emails,omitempty"`
	PhoneNumbers []PhoneNumber `json:"phoneNumbers,omitempty"`
	Groups       []GroupMember `json:"groups,omitempty"`

	EnterpriseUser *EnterpriseUserExtension `json:"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User,omitempty"`
}

func NewSCIMUserResponse(
	id, externalID, userName, displayName string,
	isActive bool,
	createdAt, lastModifiedAt time.Time,
	resourceLocation, etagVersion string,
) User {
	schemas := []string{SchemaUser}

	return User{
		Schemas:    schemas,
		ID:         id,
		ExternalID: externalID,
		Meta: Meta{
			ResourceType: "User",
			Created:      createdAt.UTC(),
			LastModified: lastModifiedAt.UTC(),
			Location:     resourceLocation,
			Version:      etagVersion,
		},
		UserName:    userName,
		DisplayName: displayName,
		Active:      isActive,
	}
}

func DummySCIMUser(userID string) User {
	if userID == "" {
		userID = "d2d46e8c-8435-4a25-a7b6-1f7c0a9e7b2f"
	}

	currentTime := time.Now().UTC()
	createdTime := currentTime.Add(-time.Hour * 24 * 30)
	lastModifiedTime := currentTime

	managerID := "a1b2c3d4-e5f6-7890-1234-567890abcdef"

	return User{
		Schemas: []string{
			SchemaUser,
			SchemaEnterpriseUser,
		},
		ID:         userID,
		ExternalID: "alice.smith.corp.id-12345",
		Meta: Meta{
			ResourceType: "User",
			Created:      createdTime,
			LastModified: lastModifiedTime,
			Location:     "https://api.example.com/scim/v2/Users/" + userID,
			Version:      "W/\"2024-07-29T14:30:00Z\"",
		},
		UserName:          "asmith",
		DisplayName:       "Alice Smith",
		NickName:          "Ally",
		ProfileURL:        "https://example.com/profiles/asmith",
		Title:             "Software Engineer",
		UserType:          "Employee",
		PreferredLanguage: "en-US",
		Locale:            "en-US",
		Timezone:          "America/Los_Angeles",
		Active:            true,
		Password:          "",

		Name: &Name{
			Formatted:       "Alice P. Smith",
			FamilyName:      "Smith",
			GivenName:       "Alice",
			MiddleName:      "P.",
			HonorificPrefix: "Ms.",
		},
		Emails: []Email{
			{
				Value:   "alice.smith@example.com",
				Type:    "work",
				Primary: true,
			},
			{
				Value:   "alice.personal@gmail.com",
				Type:    "home",
				Primary: false,
			},
		},
		PhoneNumbers: []PhoneNumber{
			{
				Value:   "+1-555-123-4567",
				Type:    "mobile",
				Primary: true,
			},
			{
				Value:   "+1-555-987-6543",
				Type:    "work",
				Primary: false,
			},
		},
		Groups: []GroupMember{
			{
				Value:   "a1b2c3d4-e5f6-7890-abcd-ef0123456789",
				Ref:     "https://api.example.com/scim/v2/Groups/a1b2c3d4-e5f6-7890-abcd-ef0123456789",
				Display: "Engineering Team",
				Type:    "Group",
			},
			{
				Value:   "fedcba98-7654-3210-fedc-ba9876543210",
				Ref:     "https://api.example.com/scim/v2/Groups/fedcba98-7654-3210-fedc-ba9876543210",
				Display: "All Employees",
				Type:    "Group",
			},
		},
		EnterpriseUser: &EnterpriseUserExtension{
			EmployeeNumber: "AES007",
			Organization:   "Example Corp",
			Department:     "Engineering",
			Division:       "Backend Services",
			CostCenter:     "CC1002",
			Manager: &Manager{
				Value:       managerID,
				Ref:         "https://api.example.com/scim/v2/Users/" + managerID,
				DisplayName: "Bob Johnson",
			},
		},
	}
}

type ListResponse struct {
	Schemas      []string `json:"schemas"`
	TotalResults int      `json:"totalResults"`
	StartIndex   int      `json:"startIndex"`
	ItemsPerPage int      `json:"itemsPerPage"`
	Resources    []User   `json:"Resources"`
}

func NewSCIMUserListResponse(
	users []User,
	totalResults int,
	startIndex int,
	itemsPerPage int,
) ListResponse {
	return ListResponse{
		Schemas:      []string{SchemaListResponse},
		TotalResults: totalResults,
		StartIndex:   startIndex,
		ItemsPerPage: itemsPerPage,
		Resources:    users,
	}
}

type UserCreateRequest struct {
	Schemas    []string `json:"schemas"`
	ExternalID string   `json:"externalId,omitempty"`

	UserName          string `json:"userName"`
	Name              *Name  `json:"name,omitempty"`
	DisplayName       string `json:"displayName,omitempty"`
	NickName          string `json:"nickName,omitempty"`
	ProfileURL        string `json:"profileUrl,omitempty"`
	Title             string `json:"title,omitempty"`
	UserType          string `json:"userType,omitempty"`
	PreferredLanguage string `json:"preferredLanguage,omitempty"`
	Locale            string `json:"locale,omitempty"`
	Timezone          string `json:"timezone,omitempty"`
	Active            bool   `json:"active,omitempty"`
	Password          string `json:"password,omitempty"`

	Emails       []Email       `json:"emails,omitempty"`
	PhoneNumbers []PhoneNumber `json:"phoneNumbers,omitempty"`
	Groups       []GroupMember `json:"groups,omitempty"`

	EnterpriseUser *EnterpriseUserExtension `json:"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User,omitempty"`
}
