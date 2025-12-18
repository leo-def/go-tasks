-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS ratings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rate DOUBLE PRECISION NOT NULL,
    participation_id UUID NOT NULL REFERENCES participations(id) ON DELETE CASCADE,
    collaborator_id UUID NOT NULL REFERENCES collaborators(id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL
);
CREATE INDEX IF NOT EXISTS idx_ratings_deleted_at ON ratings (deleted_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_ratings_deleted_at;
DROP TABLE IF EXISTS ratings;
-- +goose StatementEnd
