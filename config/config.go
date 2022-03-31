package config

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const (
	TokenHeader = "Authorization"
	Bearer      = "Bearer "
)

var (
	WorkDir string
)

func Init() {
	var err, err2, err3 error
	WorkDir = os.Getenv("TODO_GIN_WORKDIR")
	switch gin.Mode() {
	case gin.DebugMode:
		err = godotenv.Load(fmt.Sprintf("%v/config/development.env", WorkDir))
		err2 = godotenv.Load(fmt.Sprintf("%v/config/secret.env", WorkDir))
	case gin.TestMode:
		err = godotenv.Load(fmt.Sprintf("%v/config/test.env", WorkDir))
	default:
		err = godotenv.Load(fmt.Sprintf("%v/config/release.env", WorkDir))
	}

	err3 = godotenv.Load(fmt.Sprintf("%v/config/common.env", WorkDir))
	if err != nil || err2 != nil || err3 != nil {
		panic("Failed to load env file")
	}
}
