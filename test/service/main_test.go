package service_test

import (
	"github.com/golang-jwt/jwt"
	"github.com/kuritaeiji/todo-gin-back/service"
)

var (
	email       = "user0@example.com"
	password    = "Password1010"
	tokenString = "tokenstring"
	id          = 1
	claim       = &service.UserClaim{id, jwt.StandardClaims{}}
)
