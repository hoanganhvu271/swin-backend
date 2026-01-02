package config

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var FirebaseApp *firebase.App

func InitFirebase() {
	ctx := context.Background()

	// For VPS
	//credPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")

	// For local testing
	credPath := "./serviceAccount.json"

	if credPath == "" {
		log.Fatal("GOOGLE_APPLICATION_CREDENTIALS is not set")
	}

	opt := option.WithCredentialsFile(credPath)

	config := &firebase.Config{
		ProjectID:     "swin-55203",
		StorageBucket: "swin-55203.appspot.com",
	}

	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		log.Fatalf("error initializing firebase: %v", err)
	}

	FirebaseApp = app
	log.Println("Firebase initialized successfully")
}
