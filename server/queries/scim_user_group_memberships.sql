-- name: CreateUserGroupMembership :exec
INSERT INTO scim_user_group_memberships (
    user_id,
    group_id
) VALUES (
    sqlc.arg(user_id),
    sqlc.arg(group_id)
) ON CONFLICT (user_id, group_id) DO NOTHING;

-- name: GetUserGroupMemberships :many
SELECT
    user_id,
    group_id
FROM scim_user_group_memberships
WHERE user_id = sqlc.arg(user_id)
ORDER BY group_id;
