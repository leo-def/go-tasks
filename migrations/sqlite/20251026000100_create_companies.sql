-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS companies (
    id TEXT PRIMARY KEY DEFAULT (
        lower(hex(randomblob(4))) || '-' ||
        lower(hex(randomblob(2))) || '-' ||
        '4' || substr(lower(hex(randomblob(2))),2) || '-' ||
        substr('89ab',abs(random()) % 4 + 1,1) || substr(lower(hex(randomblob(2))),2) || '-' ||
        lower(hex(randomblob(6)))
    ),
    title TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL
);

CREATE INDEX IF NOT EXISTS idx_companies_deleted_at ON companies (deleted_at);
-- Simple index on title to help LIKE filters
CREATE INDEX IF NOT EXISTS idx_companies_title ON companies (title);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_companies_title;
DROP INDEX IF EXISTS idx_companies_deleted_at;
DROP TABLE IF EXISTS companies;
-- +goose StatementEnd
