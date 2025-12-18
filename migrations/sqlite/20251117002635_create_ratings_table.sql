-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS ratings (
    id TEXT PRIMARY KEY DEFAULT (
        lower(hex(randomblob(4))) || '-' ||
        lower(hex(randomblob(2))) || '-' ||
        '4' || substr(lower(hex(randomblob(2))),2) || '-' ||
        substr('89ab',abs(random()) % 4 + 1,1) || substr(lower(hex(randomblob(2))),2) || '-' ||
        lower(hex(randomblob(6)))
    ),
    rate REAL NOT NULL,
    participation_id TEXT NOT NULL REFERENCES participations(id) ON DELETE CASCADE,
    collaborator_id TEXT NOT NULL REFERENCES collaborators(id) ON DELETE RESTRICT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL
);
CREATE INDEX IF NOT EXISTS idx_ratings_deleted_at ON ratings (deleted_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_ratings_deleted_at;
DROP TABLE IF EXISTS ratings;
-- +goose StatementEnd
