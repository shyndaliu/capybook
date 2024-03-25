CREATE TABLE IF NOT EXISTS users (
id bigserial PRIMARY KEY,
username citext NOT NULL UNIQUE,
email citext NOT NULL UNIQUE,
password bytea NOT NULL,
token_hash string NOT NULL,
activated bool NOT NULL DEFAULT false
);
CREATE TABLE IF NOT EXISTS admins {
    id bigserial PRIMARY KEY,
    user_id bigint,
    FOREIGN KEY (user_id)
        REFERENCES users(id),
}