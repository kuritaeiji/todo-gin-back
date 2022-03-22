package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/kuritaeiji/todo-gin-back/config"
	"github.com/kuritaeiji/todo-gin-back/service"
)

type UserController interface {
	Create(ctx *gin.Context)
	IsUnique(ctx *gin.Context)
}

type userController struct {
	service      service.UserService
	emailService service.EmailService
}

func NewUserController() UserController {
	return &userController{service.NewUserService(), service.NewEmailService()}
}

func (c *userController) Create(ctx *gin.Context) {
	user, err := c.service.Create(ctx)
	if verr, ok := err.(validator.ValidationErrors); ok {
		ctx.JSON(400, gin.H{
			"message": verr.Error(),
		})
		return
	}
	if err == config.UniqueUserError {
		ctx.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.emailService.ActivationUserEmail(user)
	ctx.Status(200)
}

func (c *userController) IsUnique(ctx *gin.Context) {
	if c.service.IsUnique(ctx) {
		ctx.Status(200)
		return
	}

	ctx.AbortWithStatus(400)
}

// test用
func TestNewUserController(us service.UserService, es service.EmailService) UserController {
	return &userController{us, es}
}
