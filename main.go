package main

import (
	"log"

    "clinic/sqlc"
    "clinic/src/routes"
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
    AllRoutes := routes.NewAllRoutes(r,queries, redisClient)
    AllRoutes.SetUpAllRoutes()
    
    port := utils.GetPort()
    r.Run(":" + port)
}
