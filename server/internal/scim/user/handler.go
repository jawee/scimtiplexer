package user

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/jawee/scimtiplexer/internal/repository"
)

type handler struct {
	repo    repository.Querier
	service *service
}

func RegisterEndpoints(mux *http.ServeMux, repo repository.Querier) {
	h := &handler{
		repo:    repo,
		service: &service{repo: repo},
	}

	slog.Debug("Registering SCIM endpoints")
	h.registerScimEndpoints(mux)

	slog.Debug("SCIM endpoints registered")
}

var SCIM_PREFIX = "/scim/v2/"

func (s *handler) registerScimEndpoints(mux *http.ServeMux) {
	s.registerScimEndpoint(mux, "GET", "Users", http.HandlerFunc(s.handleGetUsers))
	s.registerScimEndpoint(mux, "GET", "Users/", http.HandlerFunc(s.handleGetUsers))
	s.registerScimEndpoint(mux, "POST", "Users", http.HandlerFunc(s.handlePostUsers))

	s.registerScimEndpoint(mux, "GET", "Users/{id}", http.HandlerFunc(s.handleGetUserById))
}

func (s *handler) registerScimEndpoint(mux *http.ServeMux, method, resource string, handler http.Handler) {
	mux.Handle(method+" "+SCIM_PREFIX+resource, s.ScimEndpointAuth(handler))
	mux.Handle(method+" "+SCIM_PREFIX+strings.ToLower(resource), s.ScimEndpointAuth(handler))
}

func (s *handler) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	organisationId, ok := r.Context().Value("orgid").(string)
	if !ok || organisationId == "" {
		slog.Error("Organisation ID not found in context")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	slog.Debug("handleGetUsers called for organisation", "orgid", r.Context().Value("orgid"))

	// TODO: Handle filter and attributes
	queryParams := r.URL.Query()
	slog.Debug("Query parameters", "params", queryParams)

	users, err := s.service.GetAllUsers(r.Context(), organisationId)
	if err != nil {
		slog.Error("Failed to get users", "error", err)
		if err == sql.ErrNoRows {
			slog.Info("No users found for organisation", "orgid", r.Context().Value("orgid"))
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	respUsers := make([]User, len(users))
	for i, user := range users {
		respUsers[i] = ScimUserResponse(user)
	}

	w.Header().Set("Content-Type", "application/scim+json")
	w.WriteHeader(http.StatusOK)
	userResp := NewSCIMUserListResponse(respUsers, len(respUsers), 1, len(respUsers))
	jsonOutput, _ := json.Marshal(userResp)
	w.Write(jsonOutput)
}

func (s *handler) handlePostUsers(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handlePostUsers called for organisation", "orgid", r.Context().Value("orgid"))

	var userReq UserCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
		slog.Error("Failed to decode user creation request", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	slog.Debug("User creation request", "request", userReq)

	createdUser, err := s.service.CreateUser(r.Context(), r.Context().Value("orgid").(string), userReq)
	if err != nil {
		slog.Error("Failed to create user", "error", err)
		if err == sql.ErrNoRows {
			slog.Info("User creation failed, no rows affected")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	slog.Debug("User created successfully", "userID", createdUser.ID)

	w.Header().Set("Content-Type", "application/scim+json")
	w.WriteHeader(http.StatusCreated)
	userResp := ScimUserResponse(createdUser)
	jsonOutput, _ := json.Marshal(userResp)
	w.Write(jsonOutput)
}

func (s *handler) handleGetUserById(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handleGetUserById called for organisation", "orgid", r.Context().Value("orgid"))
	requestedId := r.PathValue("id")
	slog.Debug("Requested user ID", "id", requestedId)

	user, err := s.service.GetUser(r.Context(), r.Context().Value("orgid").(string), requestedId)
	if err != nil {
		slog.Error("Failed to get user by ID", "error", err, "id", requestedId)
		if err == sql.ErrNoRows {
			slog.Info("User not found", "id", requestedId)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/scim+json")
	w.WriteHeader(http.StatusOK)
	userResp := ScimUserResponse(user)
	jsonOutput, _ := json.Marshal(userResp)
	w.Write(jsonOutput)
}

func (s *handler) ScimEndpointAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("ScimEndpointAuth called", "method", r.Method, "url", r.URL.Path)

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := s.repo.GetOrganisationTokenByToken(r.Context(), tokenStr)
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

func ScimUserResponse(user scimUserDto) User {
	schemas := []string{SchemaUser, SchemaEnterpriseUser}

	createdAt, _ := time.Parse(time.RFC3339, user.MetaCreated)
	lastModifiedAt, _ := time.Parse(time.RFC3339, user.MetaLastModified)

	usr := User{
		Schemas: schemas,
		ID:      user.ID,
		Meta: Meta{
			ResourceType: "User",
			Created:      createdAt.UTC(),
			LastModified: lastModifiedAt.UTC(),
			Location:     "https://api.example.com/scim/v2/Users/" + user.ID,
			Version:      user.MetaVersion,
		},
		UserName:          user.UserName,
		DisplayName:       user.DisplayName,
		UserType:          user.UserType,
		PreferredLanguage: user.PreferredLanguage,
		Locale:            user.Locale,
		Timezone:          user.Timezone,
		Active:            user.Active,
		ExternalID:        user.ExternalID,
		NickName:          user.NickName,
		ProfileURL:        user.ProfileUrl,
		Title:             user.Title,
		Emails:            []Email{},
		PhoneNumbers:      []PhoneNumber{},
		Groups:            []GroupMember{},
		Name: &Name{
			Formatted:       user.NameFormatted,
			FamilyName:      user.NameFamilyName,
			GivenName:       user.NameGivenName,
			MiddleName:      user.NameMiddleName,
			HonorificPrefix: user.NameHonorificPrefix,
			HonorificSuffix: user.NameHonorificSuffix,
		},
	}

	usr.EnterpriseUser = &EnterpriseUserExtension{
		EmployeeNumber: user.EmployeeNumber,
		Organization:   user.Organization,
		Department:     user.Department,
		Division:       user.Division,
		CostCenter:     user.CostCenter,
		Manager: &Manager{
			Value:       user.ManagerID,
			Ref:         "https://api.example.com/scim/v2/Users/" + user.ManagerID,
			DisplayName: "", //TODO: Fetch manager display name if available
		},
	}

	for _, email := range user.Emails {
		usr.Emails = append(usr.Emails, Email{
			Value:   email.Value,
			Display: email.DisplayName,
			Type:    email.Type,
			Primary: email.Primary,
		})
	}
	for _, phone := range user.PhoneNumbers {
		usr.PhoneNumbers = append(usr.PhoneNumbers, PhoneNumber{
			Value:   phone.Value,
			Display: phone.DisplayName,
			Type:    phone.Type,
			Primary: phone.Primary,
		})
	}

	return usr
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
