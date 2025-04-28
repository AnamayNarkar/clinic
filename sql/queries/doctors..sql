-- name: GetDoctorById :one
SELECT * FROM doctors WHERE id = $1;

-- name: GetDoctorByEmail :one
SELECT * FROM doctors WHERE email = $1;