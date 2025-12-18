-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS assignments (
    id TEXT PRIMARY KEY DEFAULT (
        lower(hex(randomblob(4))) || '-' ||
        lower(hex(randomblob(2))) || '-' ||
        '4' || substr(lower(hex(randomblob(2))),2) || '-' ||
        substr('89ab',abs(random()) % 4 + 1,1) || substr(lower(hex(randomblob(2))),2) || '-' ||
        lower(hex(randomblob(6)))
    ),
    assigned_to_id TEXT NOT NULL REFERENCES participations(id) ON DELETE CASCADE,
    assigner_id TEXT NOT NULL REFERENCES collaborators(id) ON DELETE RESTRICT,
    task_id TEXT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    assign_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    assign_end_date DATETIME NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS assignments;
-- +goose StatementEnd
