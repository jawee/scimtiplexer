-- name: CreateOrganisationToken :one
INSERT INTO organisation_tokens (id, organisation_id, token, created_by, created_on_utc, modified_on_utc, modified_by)
VALUES (sqlc.arg(id), sqlc.arg(organisationId), sqlc.arg(token), sqlc.arg(createdBy), sqlc.arg(createdOnUtc), sqlc.arg(modifiedOnUtc), sqlc.arg(modifiedBy))
RETURNING id;

-- name: GetOrganisationTokens :many
SELECT * FROM organisation_tokens
WHERE organisation_id = sqlc.arg(organisationId)
ORDER BY id DESC;
