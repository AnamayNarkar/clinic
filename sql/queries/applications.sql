-- name: GetApplication :one
SELECT * FROM applications
WHERE id = $1 LIMIT 1;

-- name: ListApplications :many
SELECT * FROM applications
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListApplicationsByPatient :many
SELECT * FROM applications
WHERE patient_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListApplicationsByDoctor :many
SELECT * FROM applications
WHERE doctor_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListApplicationsByStatus :many
SELECT * FROM applications
WHERE status = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CreateApplication :one
INSERT INTO applications (
  patient_id,
  doctor_id,
  status,
  description
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateApplicationStatus :one
UPDATE applications
SET 
  status = $2,
  updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: UpdateApplication :one
UPDATE applications
SET 
  patient_id = COALESCE($2, patient_id),
  doctor_id = COALESCE($3, doctor_id),
  status = COALESCE($4, status),
  description = COALESCE($5, description),
  updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteApplication :exec
DELETE FROM applications
WHERE id = $1;