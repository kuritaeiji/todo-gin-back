package factory

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/kuritaeiji/todo-gin-back/db"
	"github.com/kuritaeiji/todo-gin-back/dto"
	"github.com/kuritaeiji/todo-gin-back/model"
	"github.com/kuritaeiji/todo-gin-back/service"
)

type UserConfig struct {
	ID                 int
	Email              string
	Password           string
	Activated          bool
	OpenID             string
	NotUseDefaultValue bool
}

const (
	DefaultEmail    = "user@example.com"
	DefualtPassword = "Password1010"
)

var emailCount = 1

func (config *UserConfig) setDefaultValue() {
	if config.NotUseDefaultValue {
		return
	}

	if config.Email == "" {
		config.Email = fmt.Sprintf("%v%v", emailCount, DefaultEmail)
		emailCount++
	}
	if config.Password == "" {
		config.Password = DefualtPassword
	}
}

func NewDtoUser(config *UserConfig) dto.User {
	config.setDefaultValue()
	return dto.User{Email: config.Email, Password: config.Password}
}

func NewUser(config *UserConfig) model.User {
	dtoUser := NewDtoUser(config)
	var user model.User
	dtoUser.Transfer(&user)
	user.ID = config.ID
	user.Activated = config.Activated
	user.OpenID = config.OpenID

	return user
}

func CreateUser(config *UserConfig) model.User {
	user := NewUser(config)
	db.GetDB().Create(&user)
	return user
}

func CreateAccessToken(user model.User) string {
	return service.NewJWTService().CreateJWT(user, service.DayFromNowAccessToken)
}

func CreateUserClaim(user model.User) service.UserClaim {
	return service.UserClaim{
		ID:             user.ID,
		StandardClaims: jwt.StandardClaims{},
	}
}

func CreateUserRequestBody(config *UserConfig) io.Reader {
	config.setDefaultValue()
	body := map[string]string{
		"email":    config.Email,
		"password": config.Password,
	}
	bodyBytes, _ := json.Marshal(body)
	return strings.NewReader(string(bodyBytes))
}
