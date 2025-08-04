package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jawee/scimtiplexer/internal/repository"
)

type service struct {
	repo repository.Querier
}

func (s *service) GetUser(ctx context.Context, organisationId, id string) (repository.ScimUser, error) {
	user, err := s.repo.GetScimUserById(ctx, repository.GetScimUserByIdParams{
		Organisationid: organisationId,
		ID:             id,
	})

	if err != nil {
		return repository.ScimUser{}, err
	}
	return user, nil
}

func (u *UserCreateRequest) toCreateScimUserParams(organisationId string) (repository.CreateScimUserParams, error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		return repository.CreateScimUserParams{}, errors.New("failed to generate UUID for new user")
	}
	scimUser := repository.CreateScimUserParams{
		ID: 			uuid.String(), 
		OrganisationID: organisationId,

		UserName:       u.UserName,
		DisplayName:    sql.NullString{
			String: u.DisplayName,
			Valid: u.DisplayName != "",
		},
		Active:         u.Active,
	}
	return scimUser, nil
}

func (s *service) CreateUser(ctx context.Context, organisationId string, user UserCreateRequest) (repository.ScimUser, error) {
	newUser, err := user.toCreateScimUserParams(organisationId)
	if err != nil {
		return repository.ScimUser{}, fmt.Errorf("failed to convert user create request to params: %w", err)
	}
	userCreateResp, err := s.repo.CreateScimUser(ctx, newUser)
	if err != nil {
		return repository.ScimUser{}, fmt.Errorf("failed to CreateScimUser: %w", err)
	}

	if userCreateResp == "" {
		return repository.ScimUser{}, errors.New("failed to create user, no ID returned")
	}

	createdUser, err := s.repo.GetScimUserById(ctx, repository.GetScimUserByIdParams{
		Organisationid: organisationId,
		ID:             userCreateResp,
	})

	if err != nil {
		return repository.ScimUser{}, fmt.Errorf("failed to GetScimUserById: %w", err)
	}
	return createdUser, nil
}
