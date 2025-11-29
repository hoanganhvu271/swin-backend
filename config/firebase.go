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
	opt := option.WithCredentialsFile("serviceAccount.json")

	config := &firebase.Config{
		StorageBucket: "swin-55203.appspot.com",
	}

	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		log.Fatalf("error initializing firebase: %v", err)
	}

	FirebaseApp = app
}
