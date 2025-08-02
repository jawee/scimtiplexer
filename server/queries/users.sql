-- name: GetAllUsers :many
SELECT * FROM users;

-- name: RegisterUser :one
INSERT INTO users (id, username, email, password, created_by, created_on_utc, modified_on_utc, modified_by)
VALUES (sqlc.arg(id), sqlc.arg(username), sqlc.arg(email), sqlc.arg(password), sqlc.arg(createdBy), sqlc.arg(createdOnUTC), sqlc.arg(modifiedOnUTC), sqlc.arg(modifiedBy))
RETURNING id;
