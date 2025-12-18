-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS status_updates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lifecycle_id UUID NOT NULL REFERENCES lifecycles(id) ON DELETE CASCADE,
    status_before VARCHAR(64) NOT NULL,
    status_after VARCHAR(64) NOT NULL,
    update_date TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS status_updates;
-- +goose StatementEnd
