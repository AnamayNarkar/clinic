package main

import (
	"log"

	"clinic/sqlc"
	"clinic/src/routes"
	"clinic/src/security"
	"clinic/src/utils"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	r := gin.Default()
	r.Use(utils.SetupCORS())

	if err := utils.LoadEnv(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	db, err := utils.SetupDatabase()
	if err != nil {
		log.Fatalf("Error setting up database: %v", err)
	}

	defer db.Close()

	redisClient := utils.SetupRedis()
	queries := sqlc.New(db)
	securityManager := security.NewSecurityManager()
	routes.SetUpAllRoutes(r, queries, redisClient, securityManager)

	port := utils.GetPort()
	r.Run(":" + port)
}
