-- +goose Up
-- +goose StatementBegin
-- Add unique constraint to ensure one rating per collaborator per participation
CREATE UNIQUE INDEX IF NOT EXISTS unique_rating_per_collaborator_participation 
    ON ratings (participation_id, collaborator_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove the unique constraint
DROP INDEX IF EXISTS unique_rating_per_collaborator_participation;
-- +goose StatementEnd