// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: user.sql

package sqlc

import (
	"context"
	"database/sql"
	"time"
)

const createUser = `-- name: CreateUser :exec
INSERT INTO users (
    id, chat_id, created_at, updated_at, language, reminder_time 
    ) VALUES ( ?, ?, ?, ?, ?, ?)
`

type CreateUserParams struct {
	ID           string
	ChatID       int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Language     sql.NullString
	ReminderTime string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) error {
	_, err := q.db.ExecContext(ctx, createUser,
		arg.ID,
		arg.ChatID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Language,
		arg.ReminderTime,
	)
	return err
}

const getUserByChatID = `-- name: GetUserByChatID :one
SELECT id, chat_id, created_at, updated_at, language, reminder_time
    FROM users
    WHERE chat_id = ?
`

func (q *Queries) GetUserByChatID(ctx context.Context, chatID int64) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByChatID, chatID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.ChatID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Language,
		&i.ReminderTime,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, chat_id, created_at, updated_at, language, reminder_time
    FROM users
    WHERE id = ?
`

func (q *Queries) GetUserByID(ctx context.Context, id string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.ChatID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Language,
		&i.ReminderTime,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :exec
UPDATE users
    SET updated_at = ?, language = ?, reminder_time = ?
    WHERE id = ?
`

type UpdateUserParams struct {
	UpdatedAt    time.Time
	Language     sql.NullString
	ReminderTime string
	ID           string
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	_, err := q.db.ExecContext(ctx, updateUser,
		arg.UpdatedAt,
		arg.Language,
		arg.ReminderTime,
		arg.ID,
	)
	return err
}
