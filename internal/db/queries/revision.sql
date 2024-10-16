
-- name: CreateRevision :exec
INSERT 
    INTO revisions(
        id, revise_item_id, revised_at
    ) VALUES ( ?, ?, ? );

-- name: GetRevision :one
SELECT * 
    FROM revisions 
    WHERE id = ?;

-- name: GetRevisionItemRevisions :many
SELECT * 
    FROM revisions 
    WHERE revise_item_id = ?;

-- name: DeleteRevision :exec
DELETE 
    FROM revisions
    WHERE id = ?;

-- name: DeleteReviseItemRevisions :exec
DELETE 
    FROM revisions
    WHERE revise_item_id = ?;
