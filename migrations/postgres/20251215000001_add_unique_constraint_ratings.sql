-- +goose Up
-- +goose StatementBegin
-- Add unique constraint to ensure one rating per collaborator per participation
ALTER TABLE ratings ADD CONSTRAINT unique_rating_per_collaborator_participation 
    UNIQUE (participation_id, collaborator_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove the unique constraint
ALTER TABLE ratings DROP CONSTRAINT IF EXISTS unique_rating_per_collaborator_participation;
-- +goose StatementEnd