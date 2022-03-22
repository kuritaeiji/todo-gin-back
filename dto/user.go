package dto

import (
	"github.com/kuritaeiji/todo-gin-back/model"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Email    string `json:"email" binding:"required,max=100,email"`
	Password string `json:"password" binding:"required,min=8,max=50,password"`
}

func (userProxy User) Transfer(user *model.User) {
	user.Email = userProxy.Email
	digestByte, _ := bcrypt.GenerateFromPassword([]byte(userProxy.Password), bcrypt.DefaultCost)
	user.PasswordDigest = string(digestByte)
}
