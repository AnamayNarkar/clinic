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

-- NEW QUERIES FOR DOCTOR DASHBOARD
-- name: ListTodayAppointmentsByDoctor :many
SELECT * FROM appointments
WHERE doctor_id = $1 
  AND date(time_of_appointment) = CURRENT_DATE
ORDER BY time_of_appointment ASC
LIMIT $2 OFFSET $3;

-- name: CountAppointmentsByDoctorAndStatus :one
SELECT COUNT(*) FROM appointments
WHERE doctor_id = $1 AND status = $2;

-- name: ListDoctorAppointmentsStats :one
SELECT 
  COUNT(*) FILTER (WHERE status = 'scheduled') as scheduled_count,
  COUNT(*) FILTER (WHERE status = 'completed') as completed_count,
  COUNT(*) FILTER (WHERE status = 'cancelled') as cancelled_count,
  COUNT(*) FILTER (WHERE date(time_of_appointment) = CURRENT_DATE) as today_count
FROM appointments
WHERE doctor_id = $1;

-- EXISTING QUERIES CONTINUE BELOW
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
  time_of_appointment = $2,
  status = $3,
  appointment_results = $4,
  updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteAppointment :exec
DELETE FROM appointments
WHERE id = $1;