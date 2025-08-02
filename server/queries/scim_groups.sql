-- name: GetAllScimGroups :many
SELECT * FROM scim_groups
WHERE organisation_id = sqlc.arg(organisationId)
ORDER BY id;

-- name: CreateScimGroup :one
INSERT INTO scim_groups (id, display_name, external_id, meta_version, organisation_id)
VALUES (sqlc.arg(id), sqlc.arg(displayName), sqlc.arg(externalId), sqlc.arg(metaVersion), sqlc.arg(organisationId))
RETURNING id;
