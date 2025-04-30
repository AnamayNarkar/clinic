package controllers

import (
	"clinic/sqlc"
	"clinic/src/dto"
	"clinic/src/security"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type PatientController struct {
	DB          *sqlc.Queries
	RedisClient *redis.Client
}

func NewPatientController(db *sqlc.Queries, redisClient *redis.Client) *PatientController {
	return &PatientController{
		DB:          db,
		RedisClient: redisClient,
	}
}

func (pc *PatientController) CreatePatient(c *gin.Context) {
	requestBody := dto.PatientRegistrationDTO{}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}

	passwordHash, salt, err := security.HashPassword(requestBody.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": "error hashing password"})
		return
	}

	newUser := sqlc.CreatePatientParams{
		Username:     requestBody.Username,
		Email:        requestBody.Email,
		FirstName:    requestBody.FirstName,
		LastName:     requestBody.LastName,
		PasswordHash: string(passwordHash),
		Salt:         salt,
		Phone:        requestBody.Phone,
	}
	_, err2 := pc.DB.CreatePatient(c, newUser)
	if err2 != nil {
		c.JSON(500, gin.H{"error": "error creating patient"})
		return
	}

	c.JSON(201, gin.H{"message": "patient registered successfully"})
}
