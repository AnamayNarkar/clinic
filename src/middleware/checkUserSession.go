package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func GetSession(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := c.Cookie("session_id")
		if err != nil {
			c.JSON(400, gin.H{"error": "Session is required from middleware"})
			c.Abort()
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		value, err := redisClient.Get(ctx, session).Result()
		fmt.Println(value)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid session"})
			c.Abort()
			return
		}

		c.Set("session", value)
		c.Next()
	}
}
