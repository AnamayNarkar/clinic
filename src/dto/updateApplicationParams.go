package dto

import "github.com/google/uuid"

type UpdateApplicationDTO struct {
	ID                uuid.UUID `json:"id"`
	PatientID         uuid.UUID `json:"patient_id"`
	DoctorID          uuid.UUID `json:"doctor_id"`
	Status            string    `json:"status"`
	Description       string    `json:"description"`
	TimeOfAppointment string    `json:"time_of_appointment"`
}
