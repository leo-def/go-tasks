-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS lifecycles (
    id TEXT PRIMARY KEY DEFAULT (
        lower(hex(randomblob(4))) || '-' ||
        lower(hex(randomblob(2))) || '-' ||
        '4' || substr(lower(hex(randomblob(2))),2) || '-' ||
        substr('89ab',abs(random()) % 4 + 1,1) || substr(lower(hex(randomblob(2))),2) || '-' ||
        lower(hex(randomblob(6)))
    ),
    init_date DATETIME NOT NULL,
    due_date DATETIME NOT NULL,
    parent_id TEXT NULL REFERENCES lifecycles(id) ON DELETE SET NULL,
    status TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS lifecycles;
-- +goose StatementEnd
