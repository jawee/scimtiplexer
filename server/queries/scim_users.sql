-- name: GetAllScimUsers :many
SELECT * FROM scim_users
WHERE organisation_id = sqlc.arg(organisationId)
ORDER BY id;

-- name: GetScimUserById :one
SELECT * FROM scim_users
WHERE id = sqlc.arg(id)
AND organisation_id = sqlc.arg(organisationId);

-- name: CreateScimUser :one
INSERT INTO scim_users (
    id,
    external_id,
    user_name,
    display_name,
    nick_name,
    profile_url,
    title,
    user_type,
    preferred_language,
    locale,
    timezone,
    active,
    password,
    meta_resource_type,
    meta_created,
    meta_last_modified,
    meta_version,
    name_formatted,
    name_family_name,
    name_given_name,
    name_middle_name,
    name_honorific_prefix,
    name_honorific_suffix,
    employee_number,
    organization,
    department,
    division,
    cost_center,
    manager_id,
    organisation_id
) VALUES (
    sqlc.arg(id),
    sqlc.arg(external_id),
    sqlc.arg(user_name),
    sqlc.arg(display_name),
    sqlc.arg(nick_name),
    sqlc.arg(profile_url),
    sqlc.arg(title),
    sqlc.arg(user_type),
    sqlc.arg(preferred_language),
    sqlc.arg(locale),
    sqlc.arg(timezone),
    sqlc.arg(active),
    sqlc.arg(password),
    sqlc.arg(meta_resource_type),
    sqlc.arg(meta_created),
    sqlc.arg(meta_last_modified),
    sqlc.arg(meta_version),
    sqlc.arg(name_formatted),
    sqlc.arg(name_family_name),
    sqlc.arg(name_given_name),
    sqlc.arg(name_middle_name),
    sqlc.arg(name_honorific_prefix),
    sqlc.arg(name_honorific_suffix),
    sqlc.arg(employee_number),
    sqlc.arg(organization),
    sqlc.arg(department),
    sqlc.arg(division),
    sqlc.arg(cost_center),
    sqlc.arg(manager_id),
    sqlc.arg(organisation_id)
) RETURNING id;
