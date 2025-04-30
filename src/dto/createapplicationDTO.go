package dto

import (
	"github.com/google/uuid"
)

type CreateApplicationDTO struct {
	PatientID   uuid.UUID `json:"patient_id" binding:"required"`
	DoctorID    uuid.UUID `json:"doctor_id" binding:"required"`
	Status      string    `json:"status" binding:"required,oneof=PENDING APPROVED REJECTED"`
	Description string    `json:"description" binding:"omitempty"`
}
