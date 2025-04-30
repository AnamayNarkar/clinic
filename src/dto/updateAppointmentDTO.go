package dto

type UpdateAppointmentDTO struct {
	ID                 string `json:"id" `
	TimeOfAppointment  string `json:"time_of_appointment"`
	Status             string `json:"status" `
	AppointmentResults string `json:"appointment_results" `
}
