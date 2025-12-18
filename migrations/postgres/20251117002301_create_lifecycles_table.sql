-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS lifecycles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    init_date TIMESTAMPTZ NOT NULL,
    due_date TIMESTAMPTZ NOT NULL,
    parent_id UUID NULL REFERENCES lifecycles(id) ON DELETE SET NULL,
    status VARCHAR(64) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS lifecycles;
-- +goose StatementEnd
