
-- name: CreateUser :exec
INSERT INTO users (
    id, chat_id, created_at, updated_at, language, reminder_time 
    ) VALUES ( ?, ?, ?, ?, ?, ?);

-- name: GetUserByID :one
SELECT *
    FROM users
    WHERE id = ?;

-- name: GetUserByChatID :one
SELECT *
    FROM users
    WHERE chat_id = ?;

-- name: UpdateUser :exec
UPDATE users
    SET updated_at = ?, language = ?, reminder_time = ?
    WHERE id = ?;

-- name: GetUsersByReminderTime :many
SELECT *
    FROM users
    WHERE reminder_time = ?;