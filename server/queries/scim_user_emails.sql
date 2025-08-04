-- name: CreateUserEmail :exec
INSERT INTO scim_user_emails (
    id,
    user_id,
    value,
    display,
    type,
    primary_email
) VALUES (
    sqlc.arg(id),
    sqlc.arg(user_id),
    sqlc.arg(value),
    sqlc.arg(display),
    sqlc.arg(type),
    sqlc.arg(primary_email)
) ON CONFLICT (user_id, value) DO NOTHING;

-- name: GetUserEmails :many
SELECT
    *
FROM scim_user_emails
WHERE user_id = sqlc.arg(user_id)
ORDER BY value;
