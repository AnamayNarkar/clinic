-- name: GetPatient :one
SELECT * FROM patients
WHERE id = $1 LIMIT 1;

-- name: GetPatientByEmail :one
SELECT * FROM patients
WHERE email = $1 LIMIT 1;

-- name: GetPatientByUsername :one
SELECT * FROM patients
WHERE username = $1 LIMIT 1;

-- name: GetPatientByPhone :one
SELECT * FROM patients
WHERE phone = $1 LIMIT 1;

-- name: ListPatients :many
SELECT * FROM patients
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CreatePatient :one
INSERT INTO patients (
  username,
  first_name,
  last_name,
  email,
  password_hash,
  salt,
  phone
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: UpdatePatient :one
UPDATE patients
SET 
  first_name = COALESCE($2, first_name),
  last_name = COALESCE($3, last_name),
  email = COALESCE($4, email),
  password_hash = COALESCE($5, password_hash),
  salt = COALESCE($6, salt),
  phone = COALESCE($7, phone),
  updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeletePatient :exec
DELETE FROM patients
WHERE id = $1;