package controllers

import (
	"clinic/sqlc"
)

// IAppointmentController defines the interface for appointment controller operations
type IAppointmentController interface {
	CreateAppointment(body sqlc.CreateAppointmentParams) (sqlc.Appointment, error)
}
