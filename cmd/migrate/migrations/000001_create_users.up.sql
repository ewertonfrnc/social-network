CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users(
    id bigserial PRIMARY KEY,
    email citext UNIQUE NOT NULL,
    username text UNIQUE NOT NULL,
    password bytea NOT NULL,
    created_at timestamp(0) NOT NULL DEFAULT NOW()
)
