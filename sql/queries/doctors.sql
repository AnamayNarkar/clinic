-- name: GetDoctor :one
SELECT * FROM doctors
WHERE id = $1 LIMIT 1;

-- name: GetDoctorByUsername :one
SELECT * FROM doctors
WHERE username = $1 LIMIT 1;

-- name: GetDoctorByEmail :one
SELECT * FROM doctors
WHERE email = $1 LIMIT 1;

-- name: ListDoctors :many
SELECT * FROM doctors
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListDoctorsBySpecialization :many
SELECT * FROM doctors
WHERE specialization = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CreateDoctor :one
INSERT INTO doctors (
  username,
  first_name,
  last_name,
  email,
  password_hash,
  specialization,
  salt
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: UpdateDoctor :one
UPDATE doctors
SET 
  username = COALESCE($2, username),
  email = COALESCE($3, email),
  password_hash = COALESCE($4, password_hash),
  specialization = COALESCE($5, specialization),
  salt = COALESCE($6, salt),
  updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteDoctor :exec
DELETE FROM doctors
WHERE id = $1;