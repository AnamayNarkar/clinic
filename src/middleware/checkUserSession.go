package middleware

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"clinic/src/entity"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func GetSession(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := c.Cookie("session_id")
		if err != nil {
			c.JSON(401, gin.H{"error": "Session is required"})
			c.Abort()
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		value, err := redisClient.Get(ctx, session).Result()
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid session"})
			c.Abort()
			return
		}

		// Deserialize session data
		var sessionEntity entity.SessionValueEntity
		err = json.Unmarshal([]byte(value), &sessionEntity)
		if err != nil {
			log.Printf("Error deserializing session: %v", err)
			c.JSON(500, gin.H{"error": "Internal server error"})
			c.Abort()
			return
		}

		// Add session entity to context
		c.Set("session", &sessionEntity)
		c.Next()
	}
}
