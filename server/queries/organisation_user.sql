-- name: CreateOrganisationUser :exec
INSERT INTO user_organisations (user_id, organisation_id, created_on_utc, modified_on_utc)
VALUES (sqlc.arg(userId), sqlc.arg(organisationId), sqlc.arg(createdOnUTC), sqlc.arg(modifiedOnUTC));
