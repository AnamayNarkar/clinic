-- name: GetReceptionist :one
SELECT * FROM receptionists
WHERE id = $1 LIMIT 1;

-- name: GetReceptionistByUsername :one
SELECT * FROM receptionists
WHERE username = $1 LIMIT 1;

-- name: GetReceptionistByEmail :one
SELECT * FROM receptionists
WHERE email = $1 LIMIT 1;

-- name: ListReceptionists :many
SELECT * FROM receptionists
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CreateReceptionist :one
INSERT INTO receptionists (
  username,
  email,
  password_hash,
  salt
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateReceptionist :one
UPDATE receptionists
SET 
  username = COALESCE($2, username),
  email = COALESCE($3, email),
  password_hash = COALESCE($4, password_hash),
  salt = COALESCE($5, salt),
  updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteReceptionist :exec
DELETE FROM receptionists
WHERE id = $1;