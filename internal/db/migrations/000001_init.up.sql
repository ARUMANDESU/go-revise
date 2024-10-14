CREATE TABLE users (
    id TEXT PRIMARY KEY, -- UUID
    chat_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    language TEXT, -- did not yet figured it out how to store language @@TODO
    reminder_time TEXT NOT NULL -- hour:minute, e.g. 21:00 
);

CREATE TABLE revise_items (
    id TEXT PRIMARY KEY, -- UUID
    user_id TEXT NOT NULL, -- UUID
    name TEXT NOT NULL,
    description TEXT,
    tags TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    last_revised_at TIMESTAMP NOT NULL,
    next_revision_at TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE revisions (
    id TEXT PRIMARY KEY, -- UUID
    revise_item_id TEXT NOT NULL, -- UUID
    revised_ad TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (revise_item_id) REFERENCES revise_items(id)
)
