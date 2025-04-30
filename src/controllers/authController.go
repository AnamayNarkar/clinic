package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"clinic/sqlc"
	"clinic/src/dto"
	"clinic/src/entity"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	DB          *sqlc.Queries
	RedisClient *redis.Client
}

func NewAuthController(db *sqlc.Queries, redisClient *redis.Client) *AuthController {
	return &AuthController{
		DB:          db,
		RedisClient: redisClient,
	}
}

func (ac *AuthController) CreateSession(
	userID uuid.UUID,
	username string,
	role string,
	c *gin.Context,
) {
	newSessionValueEntity := entity.SessionValueEntity{
		SessionId: uuid.New().String(),
		UserId:    userID.String(),
		Username:  username,
		Role:      role,
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    newSessionValueEntity.SessionId,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	}

	http.SetCookie(c.Writer, cookie)

	sessionValue, err := json.Marshal(newSessionValueEntity)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	err = ac.RedisClient.Set(c, newSessionValueEntity.SessionId, sessionValue, 24*time.Hour).Err()
	if err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	log.Println("Session created")
}

func (ac *AuthController) Login(c *gin.Context) {
	role := c.Param("role")

	requestBody := dto.UsernameAndPasswordDTO{}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}

	var userID uuid.UUID
	var username string
	var email string
	var firstName, lastName string
	var passwordHash string
	var salt string

	switch role {
	case "doctor":
		doctor, err := ac.DB.GetDoctorByUsername(c, requestBody.Username)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(401, gin.H{"error": "invalid credentials"})
				return
			} else {
				c.JSON(500, gin.H{"error": "internal server error"})
				return
			}
		}
		userID = doctor.ID
		username = doctor.Username
		email = doctor.Email
		firstName = doctor.FirstName
		lastName = doctor.LastName
		passwordHash = doctor.PasswordHash
		salt = doctor.Salt
	case "admin":
		admin, err := ac.DB.GetAdminByUsername(c, requestBody.Username)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(401, gin.H{"error": "invalid credentials"})
				return
			} else {
				c.JSON(500, gin.H{"error": "internal server error"})
				return
			}
		}
		userID = admin.ID
		username = admin.Username
		email = admin.Email
		passwordHash = admin.PasswordHash
		salt = admin.Salt
	case "receptionist":
		receptionist, err := ac.DB.GetReceptionistByUsername(c, requestBody.Username)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(401, gin.H{"error": "invalid credentials"})
				return
			} else {
				c.JSON(500, gin.H{"error": "internal server error"})
				return
			}
		}
		userID = receptionist.ID
		username = receptionist.Username
		email = receptionist.Email
		firstName = receptionist.FirstName
		lastName = receptionist.LastName
		passwordHash = receptionist.PasswordHash
		salt = receptionist.Salt
	case "patient":
		patient, err := ac.DB.GetPatientByUsername(c, requestBody.Username)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(401, gin.H{"error": "invalid credentials"})
				return
			} else {
				c.JSON(500, gin.H{"error": "internal server error"})
				return
			}
		}
		userID = patient.ID
		username = patient.Username
		email = patient.Email
		firstName = patient.FirstName
		lastName = patient.LastName
		passwordHash = patient.PasswordHash
		salt = patient.Salt
	default:
		c.JSON(400, gin.H{"error": "invalid role"})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(requestBody.Password+salt)) != nil {
		c.JSON(401, gin.H{"error": "invalid credentials"})
		return
	}

	ac.CreateSession(userID, username, role, c)

	// Create basic user data to return
	name := username
	if firstName != "" || lastName != "" {
		name = firstName + " " + lastName
	}

	userData := gin.H{
		"id":       userID.String(),
		"username": username,
		"email":    email,
		"name":     name,
		"role":     role,
	}

	c.JSON(200, gin.H{"message": "login successful", "user": userData})
}

func (ac *AuthController) Logout(c *gin.Context) {
	cookie, err := c.Request.Cookie("session_id")
	if err != nil {
		c.JSON(400, gin.H{"error": "no session found"})
		return
	}

	err = ac.RedisClient.Del(c, cookie.Value).Err()
	if err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	cookie.Expires = time.Now().Add(-1 * time.Hour)
	http.SetCookie(c.Writer, cookie)

	c.JSON(200, gin.H{"message": "logout successful"})
}
