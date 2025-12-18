-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS status_updates (
    id TEXT PRIMARY KEY DEFAULT (
        lower(hex(randomblob(4))) || '-' ||
        lower(hex(randomblob(2))) || '-' ||
        '4' || substr(lower(hex(randomblob(2))),2) || '-' ||
        substr('89ab',abs(random()) % 4 + 1,1) || substr(lower(hex(randomblob(2))),2) || '-' ||
        lower(hex(randomblob(6)))
    ),
    lifecycle_id TEXT NOT NULL REFERENCES lifecycles(id) ON DELETE CASCADE,
    status_before TEXT NOT NULL,
    status_after TEXT NOT NULL,
    update_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS status_updates;
-- +goose StatementEnd
