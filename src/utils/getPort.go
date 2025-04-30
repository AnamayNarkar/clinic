package utils

import (
	"log"
	"os"
)

func GetPort() string {
	port := os.Getenv("PORT")
	log.Println("PORT environment variable:", port)
	if port == "" {
		log.Println("PORT environment variable not set, defaulting to 3000")
		return ":3000"
	}
	return ":" + port
}
