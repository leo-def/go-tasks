-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS companies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL
);

CREATE INDEX IF NOT EXISTS idx_companies_deleted_at ON companies (deleted_at);
-- Trigram index to optimize ILIKE '%term%' searches on title
CREATE INDEX IF NOT EXISTS idx_companies_title_trgm ON companies USING gin (title gin_trgm_ops);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_companies_title_trgm;
DROP INDEX IF EXISTS idx_companies_deleted_at;
DROP TABLE IF EXISTS companies;
-- +goose StatementEnd
