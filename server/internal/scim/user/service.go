package user

import (
	"context"

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
