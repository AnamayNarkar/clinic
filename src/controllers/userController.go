package controllers

import (
	"clinic/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type UserController struct {
	DB          *sqlc.Queries
	RedisClient *redis.Client
}

func NewUserController(db *sqlc.Queries, redisClient *redis.Client) *UserController {
	return &UserController{
		DB:          db,
		RedisClient: redisClient,
	}
}

func (uc *UserController) GetAllUsers(c *gin.Context) {
	users, err := uc.DB.GetAllUsers(c)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(200, gin.H{"users": users})
}
