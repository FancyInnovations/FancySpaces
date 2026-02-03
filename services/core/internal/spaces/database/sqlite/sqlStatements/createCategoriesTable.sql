CREATE TABLE IF NOT EXISTS space_categories (
    space_id TEXT NOT NULL,
    category TEXT NOT NULL DEFAULT 'other',
    PRIMARY KEY (space_id, category),
    FOREIGN KEY (space_id) REFERENCES spaces(id) ON DELETE CASCADE
);