CREATE TABLE IF NOT EXISTS books (
id bigserial PRIMARY KEY,
title text NOT NULL,
author text NOT NULL,
year integer NOT NULL,
description text NOT NULL,
genres text[] NOT NULL
);

