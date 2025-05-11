-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS roles
(
    id    SERIAL PRIMARY KEY,
    name  VARCHAR UNIQUE NOT NULL,
    level INT            NOT NULL
);

INSERT INTO roles (name, level)
VALUES ('user', 0);
INSERT INTO roles (name, level)
VALUES ('moderator', 10);
INSERT INTO roles (name, level)
VALUES ('admin', 25);
INSERT INTO roles (name, level)
VALUES ('creator', 100);

CREATE TABLE IF NOT EXISTS users
(
    id            SERIAL PRIMARY KEY,
    email         VARCHAR UNIQUE NOT NULL,
    password_hash VARCHAR        NOT NULL,
    role          INT DEFAULT 1 REFERENCES roles (id) ON DELETE SET DEFAULT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS roles;
-- +goose StatementEnd
