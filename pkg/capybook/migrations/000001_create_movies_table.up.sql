CREATE TABLE IF NOT EXISTS books (
id bigserial PRIMARY KEY,
title text NOT NULL,
author text NOT NULL,
year integer NOT NULL,
description text NOT NULL,
genres text[] NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
id bigserial PRIMARY KEY,
username varchar(50) NOT NULL UNIQUE,
password text NOT NULL
);

CREATE TABLE IF NOT EXISTS reviews (
id bigserial PRIMARY KEY,
created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
user_id bigserial NOT NULL,
book_id bigserial NOT NULL,
content text,
rating integer CHECK(rating>=1 and rating<=5),
FOREIGN KEY (user_id)
        REFERENCES users(id),
FOREIGN KEY (book_id)
        REFERENCES books(id)
);