-- name: CreateOrganisation :one
INSERT INTO organisations (id, name, created_by, created_on_utc, modified_on_utc, modified_by)
VALUES (sqlc.arg(id), sqlc.arg(name), sqlc.arg(createdBy), sqlc.arg(createdOnUTC), sqlc.arg(modifiedOnUTC), sqlc.arg(modifiedBy))
RETURNING id;
