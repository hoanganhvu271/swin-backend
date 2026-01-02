package main

import (
	"backend/config"
	"backend/router"
	"log"
)

func main() {
	log.Println("Starting server...")
	config.InitFirebase()
	log.Println("Firebase initialized")
	config.InitCloudinary()
	log.Println("Cloudinary initialized")

	r := router.SetupRouter()

	err := r.Run("0.0.0.0" + config.ServerPort)
	if err != nil {
		log.Println("Error starting server:", err)
		return
	}
}
