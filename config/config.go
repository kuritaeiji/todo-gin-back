package config

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const (
	TokenHeader    = "Authorization"
	Bearer         = "Bearer "
	CurrentUserKey = "currentUser"
	ListKey        = "list"
	CardKey        = "card"
	StateCookieKey = "state"
)

var (
	WorkDir          string
	CsrfCustomHeader = map[string]string{
		"key":   "X-Requested-With",
		"value": "XMLHttpRequest",
	}
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
		err2 = godotenv.Load(fmt.Sprintf("%v/config/secret.env", WorkDir))
	case gin.ReleaseMode:
		err = godotenv.Load(fmt.Sprintf("%v/config/release.env", WorkDir))
	}

	err3 = godotenv.Load(fmt.Sprintf("%v/config/common.env", WorkDir))
	if err != nil || err2 != nil || err3 != nil {
		panic("Failed to load env file")
	}
}

func MakeRandomStr(digit int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	bytes := make([]byte, digit)
	for i := range bytes {
		bytes[i] = letters[rand.Intn(len(letters))]
	}

	return string(bytes)
}
