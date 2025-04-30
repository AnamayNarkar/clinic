package controllers

import (
	"net/http"
	"strconv"

	"clinic/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AppointmentController handles appointment operations
type AppointmentController struct {
	DB *sqlc.Queries
}

// NewAppointmentController creates a new AppointmentController
func NewAppointmentController(db *sqlc.Queries) *AppointmentController {
	return &AppointmentController{DB: db}
}

// CreateAppointment creates a new appointment
func (ac *AppointmentController) CreateAppointment(c *gin.Context) {
	var req sqlc.CreateAppointmentParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	appt, err := ac.DB.CreateAppointment(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusCreated, appt)
}

// GetAppointment retrieves an appointment by ID
func (ac *AppointmentController) GetAppointment(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	appt, err := ac.DB.GetAppointment(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "appointment not found"})
		return
	}
	c.JSON(http.StatusOK, appt)
}

// ListAppointments returns paginated appointments
func (ac *AppointmentController) ListAppointments(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	appts, err := ac.DB.ListAppointments(c, sqlc.ListAppointmentsParams{Limit: int32(limit), Offset: int32(offset)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, appts)
}

// CancelAppointment deletes (cancels) an appointment
func (ac *AppointmentController) CancelAppointment(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	err = ac.DB.DeleteAppointment(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.Status(http.StatusNoContent)
}
