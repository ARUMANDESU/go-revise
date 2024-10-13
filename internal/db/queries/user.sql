-- name: SaveUser :exec
INSERT INTO users (
    id, chat_id 
) VALUES ( ?, ? )
