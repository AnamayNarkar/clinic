package controllers

import (
	"clinic/sqlc"
	"clinic/src/dto"
	"clinic/src/entity"
	"clinic/src/security"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type DoctorController struct {
	DB          *sqlc.Queries
	RedisClient *redis.Client
}

func NewDoctorController(db *sqlc.Queries, redisClient *redis.Client) *DoctorController {
	return &DoctorController{
		DB:          db,
		RedisClient: redisClient,
	}
}

func (dc *DoctorController) CreateDoctor(c *gin.Context) {
	requestBody := dto.DoctorRegistrationDTO{}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}

	passwordHash, salt, err := security.HashPassword(requestBody.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	newUser := sqlc.CreateDoctorParams{
		Username:       requestBody.Username,
		FirstName:      requestBody.FirstName,
		LastName:       requestBody.LastName,
		Email:          requestBody.Email,
		PasswordHash:   string(passwordHash),
		Specialization: requestBody.Specialization,
		Salt:           salt,
	}
	_, err2 := dc.DB.CreateDoctor(c, newUser)
	if err2 != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(201, gin.H{"message": "doctor registered successfully"})
}

func (dc *DoctorController) GetDoctorAppointments(c *gin.Context) {
	sessionInterface, exists := c.Get("session")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	session, ok := sessionInterface.(*entity.SessionValueEntity)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid session data"})
		return
	}

	docID, err := uuid.Parse(session.UserId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid doctor ID"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	appointments, err := dc.DB.ListAppointmentsByDoctor(c, sqlc.ListAppointmentsByDoctorParams{
		DoctorID: docID,
		Limit:    int32(limit),
		Offset:   int32(offset),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch appointments"})
		return
	}

	formattedAppointments := make([]gin.H, 0, len(appointments))

	for _, appt := range appointments {
		// Convert time_of_appointment to date and time
		dateTime := appt.TimeOfAppointment.Format("2006-01-02 15:04:05")
		dateParts := strings.Split(dateTime, " ")

		patientName := "Patient" // Default fallback

		formattedAppointments = append(formattedAppointments, gin.H{
			"id":          appt.ID.String(),
			"patientId":   appt.PatientID.String(),
			"patientName": patientName,
			"date":        dateParts[0],
			"time":        dateParts[1],
			"status":      appt.Status,
			"reason":      appt.Description,
			"createdAt":   appt.CreatedAt,
			"updatedAt":   appt.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, formattedAppointments)
}
