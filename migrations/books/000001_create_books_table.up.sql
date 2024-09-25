CREATE TABLE IF NOT EXISTS books (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) WITH time zone NOT NULL DEFAULT NOW(),
    title text NOT NULL,
    author text NOT NULL,
    genres text[] not null,
    version integer NOT NULL DEFAULT 1
);