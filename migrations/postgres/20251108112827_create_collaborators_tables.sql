-- +goose Up
-- +goose StatementBegin
DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'collaborator_role_enum') THEN
    CREATE TYPE collaborator_role_enum AS ENUM ('OWNER', 'MANAGER', 'OPS');
  END IF;
END $$;
 
CREATE TABLE IF NOT EXISTS collaborators (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    role collaborator_role_enum NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    activation_token TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL
);

CREATE INDEX IF NOT EXISTS idx_collaborators_deleted_at ON collaborators (deleted_at);
-- Indexes used for joins and filters on collaborator associations
CREATE INDEX IF NOT EXISTS idx_collaborators_company_id ON collaborators (company_id);
CREATE INDEX IF NOT EXISTS idx_collaborators_account_id ON collaborators (account_id);
-- Unique index for activation tokens (allow multiple NULLs via partial index)
CREATE UNIQUE INDEX IF NOT EXISTS idx_collaborators_activation_token_unique
  ON collaborators (activation_token)
  WHERE activation_token IS NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_collaborators_activation_token_unique;
DROP INDEX IF EXISTS idx_collaborators_account_id;
DROP INDEX IF EXISTS idx_collaborators_company_id;
DROP INDEX IF EXISTS idx_collaborators_deleted_at;
DROP TABLE IF EXISTS collaborators;
-- +goose StatementEnd
