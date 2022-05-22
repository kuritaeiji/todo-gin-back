package seed

import (
	"github.com/kuritaeiji/todo-gin-back/dto"
	"github.com/kuritaeiji/todo-gin-back/model"
	"github.com/kuritaeiji/todo-gin-back/repository"
)

func CreateSeedData() {
	dtoUser := dto.User{Email: "user@example.com", Password: "Password1010"}
	var user model.User
	dtoUser.Transfer(&user)
	user.Activated = true
	repository.NewUserRepository().Create(&user)
}
