package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func Init() {
	var err error
	if os.Getenv("GIN_MODE") == "release" {
		err = godotenv.Load("config/release.env")
	} else {
		err = godotenv.Load("config/development.env")
		_ = godotenv.Load("config/secret.env")
	}

	if err != nil {
		panic(fmt.Sprintf("Failed to env file\n%v", err.Error()))
	}
}
