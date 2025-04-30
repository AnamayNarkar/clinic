package controllers

import (
	"clinic/sqlc"
	"clinic/src/dto"
	"clinic/src/security"

	"github.com/gin-gonic/gin"
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
