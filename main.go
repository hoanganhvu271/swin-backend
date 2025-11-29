package main

import (
	"backend/config"
	"backend/router"
	"log"
)

func main() {
	log.Println("Starting server...")
	config.InitFirebase()

	config.InitCloudinary()

	r := router.SetupRouter()

	err := r.Run(config.ServerPort)
	if err != nil {
		log.Println("Error starting server:", err)
		return
	}
}
