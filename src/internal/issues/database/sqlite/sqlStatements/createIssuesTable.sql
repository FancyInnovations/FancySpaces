CREATE TABLE issues (
    id TEXT PRIMARY KEY,
    space TEXT NOT NULL,
    title TEXT NOT NULL CHECK (length(title) > 0 AND length(title) <= 100),
    description TEXT CHECK (length(description) <= 1000),
    type TEXT NOT NULL,
    status TEXT NOT NULL,
    priority TEXT NOT NULL,
    assignee TEXT,
    reporter TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    external_source TEXT
);