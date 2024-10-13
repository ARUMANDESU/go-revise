CREATE TABLE users (
    id TEXT PRIMARY KEY, -- UUID
    chat_id INTEGER NOT NULL
);

CREATE TABLE revise_items (
    id TEXT PRIMARY KEY, -- UUID
    user_id TEXT NOT NULL, -- UUID
    name TEXT NOT NULL,
    description TEXT,
    tags TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    last_rivised_at TIMESTAMP,
    next_revision_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE revisions (
    id TEXT PRIMARY KEY, -- UUID
    revise_item_id TEXT NOT NULL, -- UUID
    revised_ad TIMESTAMP NOT NULL,
    FOREIGN KEY (revise_item_id) REFERENCES revise_items(id)
)
