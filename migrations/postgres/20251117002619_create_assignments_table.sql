-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assigned_to_id UUID NOT NULL REFERENCES participations(id) ON DELETE CASCADE,
    assigner_id UUID NOT NULL REFERENCES collaborators(id) ON DELETE RESTRICT,
    task_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    assign_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    assign_end_date TIMESTAMPTZ NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS assignments;
-- +goose StatementEnd
