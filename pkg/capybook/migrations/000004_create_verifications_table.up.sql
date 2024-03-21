CREATE TABLE IF NOT EXISTS verifications (
code bytea PRIMARY KEY,
user_id bigint,
expiry timestamp(0) with time zone NOT NULL,
FOREIGN KEY (user_id) REFERENCES users(id) on delete CASCADE
);