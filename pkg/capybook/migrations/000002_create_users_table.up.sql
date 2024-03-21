CREATE TABLE IF NOT EXISTS users (
id bigserial PRIMARY KEY,
username citext NOT NULL UNIQUE,
email citext NOT NULL UNIQUE,
password bytea NOT NULL,
activated bool NOT NULL DEFAULT false
);