-- +goose Up
-- +goose StatementBegin

-- Create users table
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    created_by TEXT REFERENCES users(id),
    created_on_utc DATETIME NOT NULL,
    modified_on_utc DATETIME NOT NULL,
    modified_by TEXT REFERENCES users(id)
);

-- Create organisations table
CREATE TABLE organisations (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    created_by TEXT REFERENCES users(id),
    created_on_utc DATETIME NOT NULL,
    modified_on_utc DATETIME NOT NULL,
    modified_by TEXT REFERENCES users(id)
);

-- Create junction table for N:N relationship
CREATE TABLE user_organisations (
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    organisation_id TEXT NOT NULL REFERENCES organisations(id) ON DELETE CASCADE,
    created_on_utc DATETIME NOT NULL,
    modified_on_utc DATETIME NOT NULL,
    PRIMARY KEY (user_id, organisation_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS user_organisations;
DROP TABLE IF EXISTS organisations;
DROP TABLE IF EXISTS users;

-- +goose StatementEnd
