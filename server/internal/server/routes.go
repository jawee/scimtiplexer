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

// SCIM Meta object
type Meta struct {
	ResourceType string    `json:"resourceType"` // e.g., "User", "Group"
	Created      time.Time `json:"created"`
	LastModified time.Time `json:"lastModified"`
	Location     string    `json:"location"`          // URL of the resource, e.g., /Users/2819c223-7f76-453a-919d-413861904646
	Version      string    `json:"version,omitempty"` // ETag, typically optional if not using concurrency control
}

// Name complex type (for User.name)
type Name struct {
	Formatted       string `json:"formatted,omitempty"`
	FamilyName      string `json:"familyName,omitempty"`
	GivenName       string `json:"givenName,omitempty"`
	MiddleName      string `json:"middleName,omitempty"`
	HonorificPrefix string `json:"honorificPrefix,omitempty"`
	HonorificSuffix string `json:"honorificSuffix,omitempty"`
}

// Email complex type (for User.emails)
type Email struct {
	Value   string `json:"value"`
	Display string `json:"display,omitempty"`
	Type    string `json:"type,omitempty"`    // e.g., "work", "home", "other"
	Primary bool   `json:"primary,omitempty"` // Indicates if this is the primary email
}

// PhoneNumber complex type (for User.phoneNumbers)
type PhoneNumber struct {
	Value   string `json:"value"`
	Display string `json:"display,omitempty"`
	Type    string `json:"type,omitempty"`    // e.g., "work", "home", "mobile", "fax", "pager", "other"
	Primary bool   `json:"primary,omitempty"` // Indicates if this is the primary phone number
}

// GroupMember complex type (for User.groups and Group.members)
type GroupMember struct {
	Value   string `json:"value"`             // The ID of the group/user
	Ref     string `json:"$ref,omitempty"`    // The URL of the group/user, e.g., /Groups/2819c223-7f76-453a-919d-413861904646
	Display string `json:"display,omitempty"` // The display name of the group/user
	Type    string `json:"type,omitempty"`    // e.g., "User", "Group"
}

// Manager complex type for Enterprise User Extension
type Manager struct {
	Value       string `json:"value"`                 // The ID of the manager (another User resource)
	Ref         string `json:"$ref,omitempty"`        // The URL of the manager resource
	DisplayName string `json:"displayName,omitempty"` // The display name of the manager
}

// EnterpriseUserExtension represents the SCIM Enterprise User extension schema
type EnterpriseUserExtension struct {
	EmployeeNumber string   `json:"employeeNumber,omitempty"`
	Organization   string   `json:"organization,omitempty"`
	Department     string   `json:"department,omitempty"`
	Division       string   `json:"division,omitempty"`
	CostCenter     string   `json:"costCenter,omitempty"`
	Manager        *Manager `json:"manager,omitempty"` // Pointer to allow null if no manager
}

// User represents a SCIM User resource
// https://tools.ietf.org/html/rfc7643#section-4.1
type User struct {
	Schemas    []string `json:"schemas"` // MUST contain SchemaUser and other extensions used
	ID         string   `json:"id"`
	ExternalID string   `json:"externalId,omitempty"`
	Meta       Meta     `json:"meta"`

	// Core attributes
	UserName          string `json:"userName"`
	Name              *Name  `json:"name,omitempty"` // Pointer to allow null
	DisplayName       string `json:"displayName,omitempty"`
	NickName          string `json:"nickName,omitempty"`
	ProfileURL        string `json:"profileUrl,omitempty"`
	Title             string `json:"title,omitempty"`
	UserType          string `json:"userType,omitempty"`
	PreferredLanguage string `json:"preferredLanguage,omitempty"`
	Locale            string `json:"locale,omitempty"`
	Timezone          string `json:"timezone,omitempty"`
	Active            bool   `json:"active"`
	Password          string `json:"password,omitempty"` // SHOULD NOT be returned in responses except for specific scenarios

	Emails       []Email       `json:"emails,omitempty"`
	PhoneNumbers []PhoneNumber `json:"phoneNumbers,omitempty"`
	Groups       []GroupMember `json:"groups,omitempty"`
	// Other multi-valued attributes (ims, photos, addresses, entitlements, roles, x509Certificates)
	// would follow the same pattern as Emails/PhoneNumbers

	// Enterprise User Extension (note the field name matches the schema URI part after 'User:')
	// This approach directly embeds the extension. Another approach is to use a map[string]interface{}.
	// For well-known extensions, embedding is cleaner.
	EnterpriseUser *EnterpriseUserExtension `json:"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User,omitempty"`
}

// Helper function to create a basic SCIM User response object
func NewSCIMUserResponse(
	id, externalID, userName, displayName string,
	isActive bool,
	createdAt, lastModifiedAt time.Time,
	resourceLocation, etagVersion string,
	// Add other parameters as needed
) User {
	// Always include the core User schema
	schemas := []string{SchemaUser}

	// Example of conditionally adding extension schema
	// if enterpriseExtensionData is not nil { // pseudo-code
	//    schemas = append(schemas, SchemaEnterpriseUser)
	// }

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
		// Password field is typically omitted from responses.
		// It's included here only for completeness of the SCIM schema, but you should
		// handle its serialization carefully in your API handlers.
		// Password:    "", // Never return actual password
	}
}

func DummySCIMUser(userID string) User {
	if userID == "" {
		userID = "d2d46e8c-8435-4a25-a7b6-1f7c0a9e7b2f" // Default UUID if not provided
	}

	currentTime := time.Now().UTC()
	createdTime := currentTime.Add(-time.Hour * 24 * 30) // 30 days ago
	lastModifiedTime := currentTime

	managerID := "a1b2c3d4-e5f6-7890-1234-567890abcdef" // A dummy manager ID

	return User{
		Schemas: []string{
			SchemaUser,
			SchemaEnterpriseUser, // Include Enterprise User extension schema
		},
		ID:         userID,
		ExternalID: "alice.smith.corp.id-12345",
		Meta: Meta{
			ResourceType: "User",
			Created:      createdTime,
			LastModified: lastModifiedTime,
			Location:     "https://api.example.com/scim/v2/Users/" + userID, // Example location
			Version:      "W/\"2024-07-29T14:30:00Z\"",                      // Example ETag
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
		Password:          "", // Never include actual password in response

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
				Value:   "a1b2c3d4-e5f6-7890-abcd-ef0123456789", // Dummy Group ID
				Ref:     "https://api.example.com/scim/v2/Groups/a1b2c3d4-e5f6-7890-abcd-ef0123456789",
				Display: "Engineering Team",
				Type:    "Group",
			},
			{
				Value:   "fedcba98-7654-3210-fedc-ba9876543210", // Dummy Group ID
				Ref:     "https://api.example.com/scim/v2/Groups/fedcba98-7654-3210-fedc-ba9876543210",
				Display: "All Employees",
				Type:    "Group",
			},
		},
		// Enterprise User Extension
		EnterpriseUser: &EnterpriseUserExtension{
			EmployeeNumber: "AES007",
			Organization:   "Example Corp",
			Department:     "Engineering",
			Division:       "Backend Services",
			CostCenter:     "CC1002",
			Manager: &Manager{
				Value:       managerID,
				Ref:         "https://api.example.com/scim/v2/Users/" + managerID,
				DisplayName: "Bob Johnson", // Manager's display name
			},
		},
	}
}

type ListResponse struct {
	Schemas      []string `json:"schemas"`      // MUST contain SchemaListResponse
	TotalResults int      `json:"totalResults"` // The total number of results returned by the list or query operation.
	StartIndex   int      `json:"startIndex"`   // The 1-based index of the first result in the current set of results.
	ItemsPerPage int      `json:"itemsPerPage"` // The number of results in the current set of results, which MAY be less than the number of results requested.
	Resources    []User   `json:"Resources"`    // An array of SCIM resources. (Note: "Resources" capitalized as per spec)
}

// DummySCIMUser (from previous response)
// ... (Your existing DummySCIMUser function) ...

// NewSCIMUserListResponse creates a SCIM ListResponse object containing multiple User resources.
// It simulates pagination parameters for a real API response.
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
	Schemas    []string `json:"schemas"` // MUST contain SchemaUser and any extensions used
	ExternalID string   `json:"externalId,omitempty"`

	// Core attributes
	UserName          string `json:"userName"` // Required
	Name              *Name  `json:"name,omitempty"`
	DisplayName       string `json:"displayName,omitempty"`
	NickName          string `json:"nickName,omitempty"`
	ProfileURL        string `json:"profileUrl,omitempty"`
	Title             string `json:"title,omitempty"`
	UserType          string `json:"userType,omitempty"`
	PreferredLanguage string `json:"preferredLanguage,omitempty"`
	Locale            string `json:"locale,omitempty"`
	Timezone          string `json:"timezone,omitempty"`
	Active            bool   `json:"active,omitempty"`   // Defaults to true if omitted
	Password          string `json:"password,omitempty"` // Password is often provided in create requests

	Emails       []Email       `json:"emails,omitempty"`
	PhoneNumbers []PhoneNumber `json:"phoneNumbers,omitempty"`
	Groups       []GroupMember `json:"groups,omitempty"` // Groups to join on creation

	// Enterprise User Extension
	EnterpriseUser *EnterpriseUserExtension `json:"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User,omitempty"`
}
