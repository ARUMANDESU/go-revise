-- mock data

-- Insert mock users
INSERT INTO users (id, telegram_id) VALUES
('1e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e', 123456789),
('2e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e', 987654321),
('5e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e', 111111111),
('6e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e', 222222222);
-- Insert mock revise items
INSERT INTO revise_items (id, user_id, name, description, tags, iteration, created_at, updated_at, last_rivised_at, next_revision_at) VALUES
('3e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e', '1e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e', 'First Revise Item', 'Description for first revise item', 'tag1,tag2', 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('4e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e', '2e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e', 'Second Revise Item', 'Description for second revise item', 'tag3,tag4', 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('7e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e', '1e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e', 'Third Revise Item', 'Description for third revise item', 'tag5,tag6', 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('8e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e', '2e8b7e6e-8f6d-4b8e-9b8e-8f6d4b8e9b8e', 'Fourth Revise Item', 'Description for fourth revise item', 'tag7,tag8', 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);