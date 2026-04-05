-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS resumes (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL DEFAULT '',
    json_data TEXT NOT NULL,
    markdown_content TEXT NOT NULL DEFAULT '',
    template_id TEXT NOT NULL DEFAULT 'default',
    custom_css TEXT NOT NULL DEFAULT '',
    module_order TEXT NOT NULL DEFAULT '[]',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    deleted_at DATETIME
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS resumes;
-- +goose StatementEnd
