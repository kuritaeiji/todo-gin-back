package factory

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/kuritaeiji/todo-gin-back/db"
	"github.com/kuritaeiji/todo-gin-back/dto"
	"github.com/kuritaeiji/todo-gin-back/model"
	"github.com/kuritaeiji/todo-gin-back/service"
)

type UserConfig struct {
	ID        int
	Email     string
	Password  string
	Activated bool
}

const (
	DefaultEmail    = "user@example.com"
	DefualtPassword = "Password1010"
)

func NewDtoUser(config UserConfig) dto.User {
	if config.Email == "" {
		config.Email = DefaultEmail
	}
	if config.Password == "" {
		config.Password = DefualtPassword
	}

	return dto.User{Email: config.Email, Password: config.Password}
}

func NewUser(config UserConfig) model.User {
	dtoUser := NewDtoUser(config)
	var user model.User
	dtoUser.Transfer(&user)
	user.ID = config.ID
	user.Activated = config.Activated

	return user
}

func CreateUser(config UserConfig) model.User {
	user := NewUser(config)
	db.GetDB().Create(&user)
	return user
}

func CreateAccessToken(user model.User) string {
	return service.NewJWTService().CreateJWT(user, service.DayFromNowAccessToken)
}

func CreateUserRequestBody(email, password string) io.Reader {
	body := map[string]string{
		"email":    email,
		"password": password,
	}
	bodyBytes, _ := json.Marshal(body)
	return strings.NewReader(string(bodyBytes))
}
