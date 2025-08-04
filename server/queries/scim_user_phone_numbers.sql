-- name: CreateUserPhoneNumber :exec
INSERT INTO scim_user_phone_numbers (
    id,
    user_id,
    value,
    display,
    type,
    primary_phone_number
) VALUES (
    sqlc.arg(id),
    sqlc.arg(user_id),
    sqlc.arg(value),
    sqlc.arg(display),
    sqlc.arg(type),
    sqlc.arg(primary_phone_number)
    ) ON CONFLICT (user_id, value) DO NOTHING;

-- name: GetUserPhoneNumbers :many
SELECT * FROM scim_user_phone_numbers
WHERE user_id = sqlc.arg(user_id)
ORDER BY value;
