package dto

import (
	"github.com/kuritaeiji/todo-gin-back/model"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Email    string `json:"email" binding:"required,max=100,email"`
	Password string `json:"password" binding:"required,min=8,max=50,password"`
}

func (dtoUser User) Transfer(user *model.User) {
	user.Email = dtoUser.Email
	digestByte, _ := bcrypt.GenerateFromPassword([]byte(dtoUser.Password), bcrypt.DefaultCost)
	user.PasswordDigest = string(digestByte)
}
