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
	c *gin.Context,
) {
	newSessionValueEntity := entity.SessionValueEntity{
		SessionId: uuid.New().String(),
		UserId:    userID.String(),
		Username:  username,
	}

	// Set cookie
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    newSessionValueEntity.SessionId,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Path:     "/", // This is important because the cookie will be available in all routes
	}

	http.SetCookie(c.Writer, cookie)

	// Set Redis
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
	requestBody := dto.UsernameAndPasswordDto{}
	err := json.NewDecoder(c.Request.Body).Decode(&requestBody)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}

	user, err := ac.DB.GetUser(c, requestBody.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(401, gin.H{"error": "invalid credentials"})
			return
		} else {
			c.JSON(500, gin.H{"error": "internal server error"})
			return
		}
	}

	if user.Password != requestBody.Password {
		c.JSON(401, gin.H{"error": "invalid credentials"})
		return
	}

	ac.CreateSession(user.ID, user.Username, c)

	c.JSON(200, gin.H{"message": "login successful"})
}

func (ac *AuthController) Register(c *gin.Context) {
	requestBody := dto.UsernameAndPasswordDto{}
	err := json.NewDecoder(c.Request.Body).Decode(&requestBody)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}

	newUser := sqlc.CreateUserParams{
		ID:        uuid.New(),
		Username:  requestBody.Username,
		Password:  requestBody.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err2 := ac.DB.CreateUser(c, newUser)
	if err2 != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(201, gin.H{"message": "user registered successfully"})
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
