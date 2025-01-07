-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
    $1, NOW(), NOW(), $2, $3
)
RETURNING *;

-- name: GetChirps :many
SELECT * FROM chirps
ORDER BY created_at;

-- name: GetChirpById :one
SELECT * FROM chirps
WHERE id = $1;

-- name: GetUserIdByChirpId :one
SELECT user_id FROM chirps
WHERE id = $1;

-- name: DeletChirpById :exec
DELETE FROM chirps
WHERE id = $1;