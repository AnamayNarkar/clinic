package controllers

import (
	"net/http"
	"strconv"
	"time"

	"clinic/sqlc"
	"clinic/src/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ApplicationController struct {
	DB                    *sqlc.Queries
	appointmentController IAppointmentController
}

func NewApplicationController(db *sqlc.Queries, app IAppointmentController) *ApplicationController {
	return &ApplicationController{
		DB:                    db,
		appointmentController: app,
	}
}

func (ac *ApplicationController) CreateApplication(c *gin.Context) {
	var req dto.CreateApplicationDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	app, err := ac.DB.CreateApplication(c, sqlc.CreateApplicationParams{
		PatientID:   req.PatientID,
		DoctorID:    req.DoctorID,
		Status:      req.Status,
		Description: req.Description,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusCreated, app)
}

func (ac *ApplicationController) GetApplication(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	app, err := ac.DB.GetApplication(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "application not found"})
		return
	}
	c.JSON(http.StatusOK, app)
}

func (ac *ApplicationController) ListApplications(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	apps, err := ac.DB.ListApplications(c, sqlc.ListApplicationsParams{Limit: int32(limit), Offset: int32(offset)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, apps)
}

func (ac *ApplicationController) UpdateApplication(c *gin.Context) {

	var req dto.UpdateApplicationDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	app, err := ac.DB.UpdateApplication(c, sqlc.UpdateApplicationParams{
		ID:          req.ID,
		Status:      req.Status,
		Description: req.Description,
		PatientID:   req.PatientID,
		DoctorID:    req.DoctorID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error " + err.Error()})
		return
	}

	timeOfAppointment, err := time.Parse("2006-01-02T15:04:05Z", req.TimeOfAppointment)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid time format"})
		return
	}

	if app.Status == "appointment_created" {
		params := sqlc.CreateAppointmentParams{
			PatientID:         app.PatientID,
			DoctorID:          app.DoctorID,
			ApplicationID:     app.ID,
			Status:            app.Status,
			Description:       app.Description,
			TimeOfAppointment: timeOfAppointment,
		}
		_, err := ac.appointmentController.CreateAppointment(params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create appointment " + err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, app)
}
