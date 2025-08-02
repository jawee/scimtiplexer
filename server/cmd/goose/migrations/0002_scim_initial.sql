-- +goose Up
CREATE TABLE IF NOT EXISTS scim_users (
    id TEXT PRIMARY KEY,
    external_id TEXT UNIQUE,
    user_name TEXT UNIQUE NOT NULL,
    display_name TEXT,
    nick_name TEXT,
    profile_url TEXT,
    title TEXT,
    user_type TEXT,
    preferred_language TEXT,
    locale TEXT,
    timezone TEXT,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    password TEXT,
    meta_resource_type TEXT NOT NULL,
    meta_created TEXT NOT NULL,
    meta_last_modified TEXT NOT NULL,
    meta_version TEXT,
    name_formatted TEXT,
    name_family_name TEXT,
    name_given_name TEXT,
    name_middle_name TEXT,
    name_honorific_prefix TEXT,
    name_honorific_suffix TEXT,
    employee_number TEXT,
    organization TEXT,
    department TEXT,
    division TEXT,
    cost_center TEXT,
    manager_id TEXT,

    -- system fields
    organisation_id TEXT NOT NULL,
    FOREIGN KEY (organisation_id) REFERENCES organisations(id) ON DELETE CASCADE,


    -- scim foreign keys
    FOREIGN KEY (manager_id) REFERENCES scim_users(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_users_user_name ON scim_users (user_name);
CREATE INDEX IF NOT EXISTS idx_users_external_id ON scim_users (external_id);
CREATE INDEX IF NOT EXISTS idx_users_employee_number ON scim_users (employee_number);
CREATE INDEX IF NOT EXISTS idx_users_manager_id ON scim_users (manager_id);


CREATE TABLE IF NOT EXISTS scim_groups (
    id TEXT PRIMARY KEY,
    external_id TEXT UNIQUE,
    display_name TEXT UNIQUE NOT NULL,
    meta_resource_type TEXT NOT NULL,
    meta_created TEXT NOT NULL,
    meta_last_modified TEXT NOT NULL,
    meta_version TEXT,

    -- system fields
    organisation_id TEXT NOT NULL,
    FOREIGN KEY (organisation_id) REFERENCES organisations(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_groups_display_name ON scim_groups (display_name);
CREATE INDEX IF NOT EXISTS idx_groups_external_id ON scim_groups (external_id);


CREATE TABLE IF NOT EXISTS scim_user_group_memberships (
    user_id TEXT NOT NULL,
    group_id TEXT NOT NULL,
    PRIMARY KEY (user_id, group_id),
    FOREIGN KEY (user_id) REFERENCES scim_users(id) ON DELETE CASCADE,
    FOREIGN KEY (group_id) REFERENCES scim_groups(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_user_group_memberships_user_id ON scim_user_group_memberships (user_id);
CREATE INDEX IF NOT EXISTS idx_user_group_memberships_group_id ON scim_user_group_memberships (group_id);


CREATE TABLE IF NOT EXISTS scim_user_emails (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    value TEXT NOT NULL,
    display TEXT,
    type TEXT,
    primary_email BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (user_id) REFERENCES scim_users(id) ON DELETE CASCADE,
    UNIQUE(user_id, value)
);

CREATE INDEX IF NOT EXISTS idx_user_emails_user_id ON scim_user_emails (user_id);
CREATE INDEX IF NOT EXISTS idx_user_emails_value ON scim_user_emails (value);


CREATE TABLE IF NOT EXISTS scim_user_phone_numbers (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    value TEXT NOT NULL,
    display TEXT,
    type TEXT,
    primary_phone_number BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (user_id) REFERENCES scim_users(id) ON DELETE CASCADE,
    UNIQUE(user_id, value)
);

CREATE INDEX IF NOT EXISTS idx_user_phone_numbers_user_id ON scim_user_phone_numbers (user_id);
CREATE INDEX IF NOT EXISTS idx_user_phone_numbers_value ON scim_user_phone_numbers (value);


-- +goose Down
DROP TABLE IF EXISTS scim_user_phone_numbers;
DROP TABLE IF EXISTS scim_user_emails;
DROP TABLE IF EXISTS scim_user_group_memberships;
DROP TABLE IF EXISTS scim_groups;
DROP TABLE IF EXISTS scim_users;
