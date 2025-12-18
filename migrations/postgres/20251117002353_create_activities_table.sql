-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS activities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    lifecycle_id UUID NOT NULL REFERENCES lifecycles(id) ON DELETE RESTRICT,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    owner_id UUID NOT NULL REFERENCES collaborators(id) ON DELETE RESTRICT,
    created_by_id UUID NOT NULL REFERENCES collaborators(id) ON DELETE RESTRICT,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL
);
CREATE INDEX IF NOT EXISTS idx_activities_deleted_at ON activities (deleted_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_activities_deleted_at;
DROP TABLE IF EXISTS activities;
-- +goose StatementEnd
