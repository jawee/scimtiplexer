-- +goose Up
-- +goose StatementBegin

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

CREATE TABLE organisations (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    created_by TEXT REFERENCES users(id),
    created_on_utc DATETIME NOT NULL,
    modified_on_utc DATETIME NOT NULL,
    modified_by TEXT REFERENCES users(id)
);

CREATE TABLE user_organisations (
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    organisation_id TEXT NOT NULL REFERENCES organisations(id) ON DELETE CASCADE,
    created_on_utc DATETIME NOT NULL,
    modified_on_utc DATETIME NOT NULL,
    PRIMARY KEY (user_id, organisation_id)
);

CREATE TABLE organisation_tokens (
    id TEXT PRIMARY KEY,
    organisation_id TEXT NOT NULL REFERENCES organisations(id) ON DELETE CASCADE,
    token TEXT NOT NULL UNIQUE,
    created_by TEXT NOT NULL REFERENCES users(id),
    created_on_utc DATETIME NOT NULL,
    modified_on_utc DATETIME NOT NULL,
    modified_by TEXT REFERENCES users(id)
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS organisation_tokens;
DROP TABLE IF EXISTS user_organisations;
DROP TABLE IF EXISTS organisations;
DROP TABLE IF EXISTS users;

-- +goose StatementEnd
