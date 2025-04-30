package controllers

import (
	"clinic/src/dto"
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"clinic/sqlc"
	"clinic/src/entity"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AppointmentController struct {
	DB *sqlc.Queries
}

func NewAppointmentController(db *sqlc.Queries) *AppointmentController {
	return &AppointmentController{DB: db}
}

func (ac *AppointmentController) CreateAppointment(body sqlc.CreateAppointmentParams) (sqlc.Appointment, error) {

	appt, err := ac.DB.CreateAppointment(context.Background(), body)
	if err != nil {
		return sqlc.Appointment{}, err
	}
	return appt, nil

}

func (ac *AppointmentController) GetUpcomingAppointments(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	// Get upcoming appointments
	appointments, err := ac.DB.ListUpcomingAppointments(c, sqlc.ListUpcomingAppointmentsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve upcoming appointments"})
		return
	}

	c.JSON(http.StatusOK, appointments)
}

func (ac *AppointmentController) UpdateAppointmentByDoctor(c *gin.Context) {
	session, exists := c.Get("session")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No session found"})
		return
	}

	sve, ok := session.(*entity.SessionValueEntity)
	if !ok || sve.Role != "doctor" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only doctors can update their appointments"})
		return
	}

	doctorID, err := uuid.Parse(sve.UserId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor ID in session"})
		return
	}

	appointmentIDParam := c.Param("id")
	appointmentID, err := uuid.Parse(appointmentIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid appointment ID"})
		return
	}

	appointment, err := ac.DB.GetAppointment(c, appointmentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve appointment"})
		}
		return
	}

	// Check if the logged in doctor is the assigned doctor for this appointment
	if appointment.DoctorID != doctorID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update appointments assigned to you"})
		return
	}

	// Parse the request body
	var req dto.UpdateAppointmentDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	appointmentTime, err := time.Parse(time.RFC3339, req.TimeOfAppointment)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time format. Use RFC3339 format (e.g., 2006-01-02T15:04:05Z07:00)"})
		return
	}

	// Prepare the update parameters
	params := sqlc.UpdateAppointmentParams{
		ID:                 appointmentID,
		TimeOfAppointment:  appointmentTime,
		Status:             req.Status,
		AppointmentResults: req.AppointmentResults,
	}

	// Update the appointment
	updatedAppointment, err := ac.DB.UpdateAppointment(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update appointment"})
		return
	}

	c.JSON(http.StatusOK, updatedAppointment)
}
