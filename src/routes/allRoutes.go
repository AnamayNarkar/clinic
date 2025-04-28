package routes

import (
	"clinic/sqlc"
	"clinic/src/controllers"
	"clinic/src/middleware"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type AllRoutes struct {
	ginEngine   *gin.Engine
	db          *sqlc.Queries
	redisClient *redis.Client
}

func NewAllRoutes(r *gin.Engine, db *sqlc.Queries, redisClient *redis.Client) *AllRoutes {
	return &AllRoutes{
		ginEngine:   r,
		db:          db,
		redisClient: redisClient,
	}
}

func (ar *AllRoutes) SetUpAllRoutes() {
	userController := controllers.NewUserController(ar.db, ar.redisClient)

	userGroup := ar.ginEngine.Group("/api/user")
	userGroup.Use(middleware.GetSession(ar.redisClient))
	userGroup.GET("/getAll", userController.GetAllUsers)

	authController := controllers.NewAuthController(ar.db, ar.redisClient)

	authGroup := ar.ginEngine.Group("/api/auth")
	authGroup.POST("/login", authController.Login)
	authGroup.POST("/register", authController.Register)
	authGroup.POST("/logout", authController.Logout)
}
