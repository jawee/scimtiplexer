package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jawee/scimtiplexer/internal/database"
	"github.com/jawee/scimtiplexer/internal/repository"
)

func main() {
	fmt.Printf("Seeding\n")
	dbService := database.New()
	defer dbService.Close()

	repo := dbService.GetRepository()

	userId, _ := uuid.NewV7()
	ctx := context.Background()
	repo.RegisterUser(ctx, repository.RegisterUserParams{
		ID:            userId.String(),
		Username:      "testuser",
		Email:         "test@test.se",
		Password:      "test1234",
		Createdonutc:  time.Now().UTC(),
		Modifiedonutc: time.Now().UTC(),
	})

	orgId, _ := uuid.NewV7()
	repo.CreateOrganisation(ctx, repository.CreateOrganisationParams{
		ID:            orgId.String(),
		Name:          "Test Organisation",
		Createdonutc:  time.Now().UTC(),
		Modifiedonutc: time.Now().UTC(),
	})

	repo.CreateOrganisationUser(ctx, repository.CreateOrganisationUserParams{
		Organisationid: orgId.String(),
		Userid:         userId.String(),
		Createdonutc:   time.Now().UTC(),
		Modifiedonutc:  time.Now().UTC(),
	})

	fmt.Printf("Seeding completed\n")
}
