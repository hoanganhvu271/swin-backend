package config

import (
	"log"

	"github.com/cloudinary/cloudinary-go/v2"
)

var CLD *cloudinary.Cloudinary

func InitCloudinary() {
	cld, err := cloudinary.NewFromParams(
		"key",
		"key",
		"6Dqn-key",
	)
	if err != nil {
		log.Fatalf("Failed to init Cloudinary: %v", err)
	}
	CLD = cld
}
