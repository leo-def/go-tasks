-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS activities (
    id TEXT PRIMARY KEY DEFAULT (
        lower(hex(randomblob(4))) || '-' ||
        lower(hex(randomblob(2))) || '-' ||
        '4' || substr(lower(hex(randomblob(2))),2) || '-' ||
        substr('89ab',abs(random()) % 4 + 1,1) || substr(lower(hex(randomblob(2))),2) || '-' ||
        lower(hex(randomblob(6)))
    ),
    title TEXT NOT NULL,
    lifecycle_id TEXT NOT NULL REFERENCES lifecycles(id) ON DELETE RESTRICT,
    company_id TEXT NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    owner_id TEXT NOT NULL REFERENCES collaborators(id) ON DELETE RESTRICT,
    created_by_id TEXT NOT NULL REFERENCES collaborators(id) ON DELETE RESTRICT,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL
);
CREATE INDEX IF NOT EXISTS idx_activities_deleted_at ON activities (deleted_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_activities_deleted_at;
DROP TABLE IF EXISTS activities;
-- +goose StatementEnd
