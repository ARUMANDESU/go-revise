INSERT INTO users (id, chat_id, language, reminder_time)
VALUES 
    ('e471de92-5652-46b4-94e9-5ad1766874f7', 123456789, 'en', '21:00'),
    ('b0fca268-3772-407e-b446-b41ba44bf33d', 987654321, 'fr', '07:30'),
    ('50fcccfc-067a-4757-b508-c08a4a33fb06', 135792468, 'es', '18:45');

INSERT INTO revise_items (id, user_id, name, description, tags, last_revised_at, next_revision_at)
VALUES 
    ('d7accc08-981f-4aa7-8477-b1840b9a2611', 'e471de92-5652-46b4-94e9-5ad1766874f7', 'Math Basics', 'Basic math revision items', 'math, basics', '2024-10-01 08:00:00', '2024-11-01 08:00:00'),
    ('e6ff2ac2-f4d1-4fcf-ae41-5509291dd799', 'e471de92-5652-46b4-94e9-5ad1766874f7', 'Physics Fundamentals', 'Introductory physics concepts', 'physics, fundamentals', '2024-10-02 09:00:00', '2024-11-02 09:00:00'),
    ('50fcccfc-067a-4757-b508-c08a4a33fb06', 'b0fca268-3772-407e-b446-b41ba44bf33d', 'French Grammar', 'Advanced grammar rules', 'french, grammar', '2024-10-05 10:00:00', '2024-11-05 10:00:00'),
    ('b0fca268-3772-407e-b446-b41ba44bf33d', '50fcccfc-067a-4757-b508-c08a4a33fb06', 'History Overview', 'World history basics', 'history, world', '2024-10-08 11:00:00', '2024-11-08 11:00:00');

INSERT INTO revisions (id, revise_item_id, revised_at)
VALUES 
    ('d7accc08-981f-4aa7-8477-b1840b9a2611', 'd7accc08-981f-4aa7-8477-b1840b9a2611', '2024-10-15 08:00:00'),
    ('e6ff2ac2-f4d1-4fcf-ae41-5509291dd799', 'd7accc08-981f-4aa7-8477-b1840b9a2611', '2024-10-20 08:00:00'),
    ('50fcccfc-067a-4757-b508-c08a4a33fb06', 'e6ff2ac2-f4d1-4fcf-ae41-5509291dd799', '2024-10-25 09:00:00'),
    ('b0fca268-3772-407e-b446-b41ba44bf33d', '50fcccfc-067a-4757-b508-c08a4a33fb06', '2024-10-28 10:00:00'),
    ('e471de92-5652-46b4-94e9-5ad1766874f7', 'b0fca268-3772-407e-b446-b41ba44bf33d', '2024-10-30 11:00:00');
