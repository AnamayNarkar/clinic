package utils

import (
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
)

func SetupCORS() gin.HandlerFunc {
    corsConfig := cors.New(cors.Config{
        AllowOrigins:     []string{"http://127.0.0.1:5500"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Accept", "Authorization", "Origin", "Content-Type", "X-CSRF-Token"},
        ExposeHeaders:    []string{"Link"},
        AllowCredentials: true,
        MaxAge:           300,
    })
    return corsConfig
}
