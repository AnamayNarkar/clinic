-- name: GetAppointment :one
SELECT * FROM appointments
WHERE id = $1 LIMIT 1;

-- name: ListAppointments :many
SELECT * FROM appointments
ORDER BY time_of_appointment ASC
LIMIT $1 OFFSET $2;

-- name: ListAppointmentsByPatient :many
SELECT * FROM appointments
WHERE patient_id = $1
ORDER BY time_of_appointment ASC
LIMIT $2 OFFSET $3;

-- name: ListAppointmentsByDoctor :many
SELECT * FROM appointments
WHERE doctor_id = $1
ORDER BY time_of_appointment ASC
LIMIT $2 OFFSET $3;

-- name: ListAppointmentsByStatus :many
SELECT * FROM appointments
WHERE status = $1
ORDER BY time_of_appointment ASC
LIMIT $2 OFFSET $3;

-- name: ListUpcomingAppointments :many
SELECT * FROM appointments
WHERE time_of_appointment > NOW()
ORDER BY time_of_appointment ASC
LIMIT $1 OFFSET $2;

-- name: ListUpcomingAppointmentsByDoctor :many
SELECT * FROM appointments
WHERE doctor_id = $1 AND time_of_appointment > NOW()
ORDER BY time_of_appointment ASC
LIMIT $2 OFFSET $3;

-- name: ListUpcomingAppointmentsByPatient :many
SELECT * FROM appointments
WHERE patient_id = $1 AND time_of_appointment > NOW()
ORDER BY time_of_appointment ASC
LIMIT $2 OFFSET $3;

-- name: CreateAppointment :one
INSERT INTO appointments (
  patient_id,
  doctor_id,
  application_id,
  time_of_appointment,
  status,
  description,
  appointment_results
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: UpdateAppointmentStatus :one
UPDATE appointments
SET 
  status = $2,
  updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: UpdateAppointmentResults :one
UPDATE appointments
SET 
  appointment_results = $2,
  status = 'completed',
  updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: UpdateAppointment :one
UPDATE appointments
SET 
  time_of_appointment = COALESCE($2, time_of_appointment),
  status = COALESCE($3, status),
  description = COALESCE($4, description),
  appointment_results = COALESCE($5, appointment_results),
  updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteAppointment :exec
DELETE FROM appointments
WHERE id = $1;