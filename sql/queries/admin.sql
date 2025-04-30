-- name: GetAdminByUsername :one
SELECT * FROM admins
WHERE username = $1 LIMIT 1;

-- name: GetAdmin :one
SELECT * FROM admins
WHERE id = $1 LIMIT 1;