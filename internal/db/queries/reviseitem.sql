
-- name: SaveReviseItem :exec
INSERT 
    INTO revise_items (
        id, user_id, name, description, tags,
        created_at, updated_at, last_revised_at, next_revision_at
    ) VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ? );

-- name: GetReviseItem :one
SELECT * 
    FROM revise_items
    WHERE id = ?;


-- name: GetUserReviseItems :many
SELECT * 
    FROM revise_items
    WHERE user_id = ?;

-- name: UpdateReviceItem :exec
UPDATE revise_items
    SET 
        name = ?, description = ?, tags = ?, created_at = ?, 
        updated_at = ?, last_revised_at = ?, next_revision_at = ?
    WHERE id = ?;

-- name: MarkReviseItemDeleted :exec
UPDATE revise_items
    SET deleted_at = ?
    WHERE id = ?;

-- name: DeleteReviseItem :exec
DELETE 
    FROM revise_items
    WHERE id = ?;
