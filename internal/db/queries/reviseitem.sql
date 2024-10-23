
-- name: SaveReviseItem :exec
INSERT 
    INTO revise_items (
        id, user_id, name, description, tags,
        created_at, updated_at, last_revised_at, next_revision_at
    ) VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ? );

-- name: GetReviseItem :one
SELECT * 
    FROM revise_items
    WHERE id = ? AND deleted_at IS NULL;


-- name: GetUserReviseItems :many
SELECT * 
    FROM revise_items
    WHERE user_id = ? AND deleted_at IS NULL;

-- name: UpdateReviseItem :exec
UPDATE revise_items
    SET 
        name = ?, description = ?, tags = ?, created_at = ?, 
        updated_at = ?, last_revised_at = ?, next_revision_at = ?
    WHERE id = ? AND deleted_at IS NULL;

-- name: MarkReviseItemDeleted :exec
UPDATE revise_items
    SET deleted_at = ?
    WHERE id = ?;

-- name: DeleteReviseItem :exec
DELETE 
    FROM revise_items
    WHERE id = ?;

-- name: ListUserReviseItems :many
SELECT COUNT(*), *
    FROM revise_items
    WHERE user_id = ? AND deleted_at IS NULL
    ORDER BY created_at DESC
    LIMIT ? OFFSET ?;

-- name: GetUserReviseItemsByTime :many
SELECT *
    FROM revise_items
    WHERE user_id = ? AND deleted_at IS NULL AND next_revision_at <= ?;