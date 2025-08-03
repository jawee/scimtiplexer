package user

import (
	"context"

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

func (s *service) CreateUser(ctx context.Context, organisationId string, user UserCreateRequest) (repository.ScimUser, error) {
	id, err := uuid.NewV7()
	newUser := repository.CreateScimUserParams{
		ID: id.String(),
		OrganisationID: organisationId,
	}

	userCreateResp, err := s.repo.CreateScimUser(ctx, newUser)
	if err != nil {
		return repository.ScimUser{}, err
	}

	createdUser, err := s.repo.GetScimUserById(ctx, repository.GetScimUserByIdParams{
		Organisationid: organisationId,
		ID:             userCreateResp.ID,
	})
	if err != nil {
		return repository.ScimUser{}, err
	}
	return createdUser, nil
}
