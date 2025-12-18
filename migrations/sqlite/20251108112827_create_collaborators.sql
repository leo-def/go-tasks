-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS collaborators (
    id TEXT PRIMARY KEY DEFAULT (
        lower(hex(randomblob(4))) || '-' ||
        lower(hex(randomblob(2))) || '-' ||
        '4' || substr(lower(hex(randomblob(2))),2) || '-' ||
        substr('89ab',abs(random()) % 4 + 1,1) || substr(lower(hex(randomblob(2))),2) || '-' ||
        lower(hex(randomblob(6)))
    ),
    company_id TEXT NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    account_id TEXT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    role TEXT NOT NULL CHECK (role IN ('OWNER','MANAGER','OPS')),
    active BOOLEAN NOT NULL DEFAULT TRUE,
    activation_token TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL
);

CREATE INDEX IF NOT EXISTS idx_collaborators_deleted_at ON collaborators (deleted_at);
-- Indexes for association fields used in joins/filters
CREATE INDEX IF NOT EXISTS idx_collaborators_company_id ON collaborators (company_id);
CREATE INDEX IF NOT EXISTS idx_collaborators_account_id ON collaborators (account_id);
-- Unique index for activation tokens
CREATE UNIQUE INDEX IF NOT EXISTS idx_collaborators_activation_token_unique ON collaborators (activation_token);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_collaborators_activation_token_unique;
DROP INDEX IF EXISTS idx_collaborators_account_id;
DROP INDEX IF EXISTS idx_collaborators_company_id;
DROP INDEX IF EXISTS idx_collaborators_deleted_at;
DROP TABLE IF EXISTS collaborators;
-- +goose StatementEnd
