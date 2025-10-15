CREATE TABLE IF NOT EXISTS spaces (
    id TEXT PRIMARY KEY,
    slug TEXT UNIQUE NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    icon_url TEXT,
    status TEXT NOT NULL DEFAULT 'draft',
    created_at DATETIME NOT NULL
);