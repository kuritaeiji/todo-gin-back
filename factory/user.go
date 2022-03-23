package factory

import (
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

func NewUser(config UserConfig) model.User {
	if config.Email == "" {
		config.Email = "user@example.com"
	}
	if config.Password == "" {
		config.Password = "Password1010"
	}

	dtoUser := dto.User{Email: config.Email, Password: config.Password}
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
