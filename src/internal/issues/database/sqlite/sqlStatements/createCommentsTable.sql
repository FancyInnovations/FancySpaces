CREATE TABLE issue_comments (
    id TEXT PRIMARY KEY,
    issue TEXT NOT NULL,
    author TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (issue) REFERENCES issues(id) ON DELETE CASCADE
);