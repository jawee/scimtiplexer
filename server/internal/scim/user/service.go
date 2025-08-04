package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jawee/scimtiplexer/internal/repository"
)

type service struct {
	repo repository.Querier
}

func (s *service) GetUser(ctx context.Context, organisationId, id string) (scimUserDto, error) {
	user, err := s.repo.GetScimUserById(ctx, repository.GetScimUserByIdParams{
		Organisationid: organisationId,
		ID:             id,
	})

	if err != nil {
		return scimUserDto{}, err
	}
	userEmails, err := s.repo.GetUserEmails(ctx, id)
	if err != nil {
		slog.Error("failed to GetUserEmails", "error", err, "userId", id)
	}
	userPhoneNumbers, err := s.repo.GetUserPhoneNumbers(ctx, id)
	if err != nil {
		slog.Error("failed to GetUserPhoneNumbers", "error", err, "userId", id)
	}

	dto := newScimUserDto(user, userEmails, userPhoneNumbers)
	return dto, nil
}

func (u *UserCreateRequest) toCreateScimUserParams(organisationId string) (repository.CreateScimUserParams, error) {
	userId, err := uuid.NewV7()
	if err != nil {
		return repository.CreateScimUserParams{}, errors.New("failed to generate UUID for new user")
	}

	scimUser := repository.CreateScimUserParams{
		ID:             userId.String(),
		OrganisationID: organisationId,

		UserName: u.UserName,
		DisplayName: sql.NullString{
			String: u.DisplayName,
			Valid:  u.DisplayName != "",
		},
		Active: u.Active,
	}

	return scimUser, nil
}

type scimUserEmailsDto struct {
	ID          string
	DisplayName string
	Type        string
	Value       string
	Primary     bool
}

type scimUserPhoneNumbersDto struct {
	ID          string
	DisplayName string
	Type        string
	Value       string
	Primary     bool
}
type scimUserDto struct {
	ID                  string
	DisplayName         string
	UserName            string
	Active              bool
	Emails              []scimUserEmailsDto
	PhoneNumbers        []scimUserPhoneNumbersDto
	ExternalID          string
	NickName            string
	ProfileUrl          string
	Title               string
	UserType            string
	PreferredLanguage   string
	Locale              string
	Timezone            string
	MetaResourceType    string
	MetaCreated         string
	MetaLastModified    string
	MetaVersion         string
	NameFormatted       string
	NameFamilyName      string
	NameGivenName       string
	NameMiddleName      string
	NameHonorificPrefix string
	NameHonorificSuffix string
	EmployeeNumber      string
	Organization        string
	Department          string
	Division            string
	CostCenter          string
	ManagerID           string
	OrganisationID      string
}

func newScimUserDto(user repository.ScimUser, emails []repository.ScimUserEmail, phoneNumbers []repository.ScimUserPhoneNumber) scimUserDto {
	dto := scimUserDto{
		ID:                  user.ID,
		DisplayName:         user.DisplayName.String,
		UserName:            user.UserName,
		Active:              user.Active,
		ExternalID:          user.ExternalID.String,
		NickName:            user.NickName.String,
		ProfileUrl:          user.ProfileUrl.String,
		Title:               user.Title.String,
		UserType:            user.UserType.String,
		PreferredLanguage:   user.PreferredLanguage.String,
		Locale:              user.Locale.String,
		Timezone:            user.Timezone.String,
		MetaResourceType:    user.MetaResourceType,
		MetaCreated:         user.MetaCreated,
		MetaLastModified:    user.MetaLastModified,
		MetaVersion:         user.MetaVersion.String,
		NameFormatted:       user.NameFormatted.String,
		NameFamilyName:      user.NameFamilyName.String,
		NameGivenName:       user.NameGivenName.String,
		NameMiddleName:      user.NameMiddleName.String,
		NameHonorificPrefix: user.NameHonorificPrefix.String,
		NameHonorificSuffix: user.NameHonorificSuffix.String,
		EmployeeNumber:      user.EmployeeNumber.String,
		Organization:        user.Organization.String,
		Department:          user.Department.String,
		Division:            user.Division.String,
		CostCenter:          user.CostCenter.String,
		ManagerID:           user.ManagerID.String,
	}

	for _, email := range emails {
		dto.Emails = append(dto.Emails, scimUserEmailsDto{
			ID:          email.ID,
			DisplayName: email.Display.String,
			Type:        email.Type.String,
			Value:       email.Value,
			Primary:     email.PrimaryEmail.Bool,
		})
	}

	for _, phone := range phoneNumbers {
		dto.PhoneNumbers = append(dto.PhoneNumbers, scimUserPhoneNumbersDto{
			ID:          phone.ID,
			DisplayName: phone.Display.String,
			Type:        phone.Type.String,
			Value:       phone.Value,
			Primary:     phone.PrimaryPhoneNumber.Bool,
		})
	}

	return dto
}

func (s *service) CreateUser(ctx context.Context, organisationId string, user UserCreateRequest) (scimUserDto, error) {
	newUser, err := user.toCreateScimUserParams(organisationId)
	if err != nil {
		return scimUserDto{}, fmt.Errorf("failed to convert user create request to params: %w", err)
	}
	userCreateResp, err := s.repo.CreateScimUser(ctx, newUser)
	if err != nil {
		return scimUserDto{}, fmt.Errorf("failed to CreateScimUser: %w", err)
	}

	if userCreateResp == "" {
		return scimUserDto{}, errors.New("failed to create user, no ID returned")
	}

	for _, email := range user.Emails {
		emailId, err := uuid.NewV7()
		if err != nil {
			slog.Error("failed to generate UUID for email", "error", err, "email", email, "userId", userCreateResp)
			continue
		}
		param := repository.CreateUserEmailParams{
			ID:     emailId.String(),
			UserID: userCreateResp,
			Display: sql.NullString{
				String: email.Display,
				Valid:  email.Display != "",
			},
			Type: sql.NullString{
				String: email.Type,
				Valid:  email.Type != "",
			},
			Value: email.Value,
			PrimaryEmail: sql.NullBool{
				Bool:  email.Primary,
				Valid: true,
			},
		}

		err = s.repo.CreateUserEmail(ctx, param)
		if err != nil {
			slog.Error("failed to create user email", "error", err, "email", email, "userId", userCreateResp)
		}
	}

	for _, phone := range user.PhoneNumbers {
		phoneId, err := uuid.NewV7()
		if err != nil {
			slog.Error("failed to generate UUID for phone", "error", err, "phone", phone, "userId", userCreateResp)
			continue
		}
		param := repository.CreateUserPhoneNumberParams{
			ID:     phoneId.String(),
			UserID: userCreateResp,
			Display: sql.NullString{
				String: phone.Display,
				Valid:  phone.Display != "",
			},
			Type: sql.NullString{
				String: phone.Type,
				Valid:  phone.Type != "",
			},
			Value: phone.Value,
			PrimaryPhoneNumber: sql.NullBool{
				Bool:  phone.Primary,
				Valid: true,
			},
		}
		err = s.repo.CreateUserPhoneNumber(ctx, param)
		if err != nil {
			slog.Error("failed to create user phone", "error", err, "phone",
				phone, "userId", userCreateResp)
		}
	}

	createdUser, err := s.repo.GetScimUserById(ctx, repository.GetScimUserByIdParams{
		Organisationid: organisationId,
		ID:             userCreateResp,
	})

	if err != nil {
		return scimUserDto{}, fmt.Errorf("failed to GetScimUserById: %w", err)
	}

	userEmails, err := s.repo.GetUserEmails(ctx, userCreateResp)
	if err != nil {
		slog.Error("failed to GetUserEmails", "error", err, "userId", userCreateResp)
	}

	userPhoneNumbers, err := s.repo.GetUserPhoneNumbers(ctx, userCreateResp)
	if err != nil {
		slog.Error("failed to GetUserPhoneNumbers", "error", err, "userId", userCreateResp)
	}

	userDto := newScimUserDto(createdUser, userEmails, userPhoneNumbers)

	return userDto, nil
}
