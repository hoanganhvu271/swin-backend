package config

import (
	"log"

	"github.com/cloudinary/cloudinary-go/v2"
)

var CLD *cloudinary.Cloudinary

func InitCloudinary() {
	cld, err := cloudinary.NewFromParams(
		"dyc1c2elf",
		"843265166772128",
		"6Dqn-lTuycx17IJsjAMnpAcGP9s",
	)
	if err != nil {
		log.Fatalf("Failed to init Cloudinary: %v", err)
	}
	CLD = cld
}
