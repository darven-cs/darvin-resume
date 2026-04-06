-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS snapshots (
    id TEXT PRIMARY KEY,
    resume_id TEXT NOT NULL,
    label TEXT NOT NULL DEFAULT '',
    note TEXT NOT NULL DEFAULT '',
    trigger_type TEXT NOT NULL DEFAULT 'manual',
    json_data TEXT NOT NULL,
    markdown_content TEXT NOT NULL,
    template_id TEXT NOT NULL DEFAULT 'minimal',
    custom_css TEXT NOT NULL DEFAULT '',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (resume_id) REFERENCES resumes(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_snapshots_resume_id ON snapshots(resume_id);
CREATE INDEX IF NOT EXISTS idx_snapshots_created_at ON snapshots(created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS snapshots;
-- +goose StatementEnd
