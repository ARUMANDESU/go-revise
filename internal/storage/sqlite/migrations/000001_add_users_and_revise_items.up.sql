CREATE TABLE users (
    id TEXT PRIMARY KEY, -- UUID
    telegram_id INTEGER NOT NULL,
);

CREATE TABLE revise_items (
    id TEXT PRIMARY KEY, -- UUID
    user_id TEXT NOT NULL, -- UUID
    name TEXT NOT NULL,
    description TEXT,
    tags TEXT,
    iteration INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_rivised_at TIMESTAMP,
    next_revision_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);