-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS ai_messages (
    id TEXT PRIMARY KEY,
    resume_id TEXT NOT NULL,
    role TEXT NOT NULL CHECK(role IN ('user','assistant')),
    content TEXT NOT NULL,
    quoted_text TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (resume_id) REFERENCES resumes(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS ai_messages;
-- +goose StatementEnd
