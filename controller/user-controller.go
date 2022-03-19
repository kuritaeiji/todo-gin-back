package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/kuritaeiji/todo-gin-back/service"
)

type UserController interface {
	Create(ctx *gin.Context)
	IsUniqueEmail(ctx *gin.Context)
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
		ctx.AbortWithStatusJSON(400, gin.H{
			"message": verr.Error(),
		})
		return
	}
	c.emailService.ActivationUserEmail(user)
	ctx.JSON(200, user)
}

func (c *userController) IsUniqueEmail(ctx *gin.Context) {
	if c.service.IsUniqueEmail(ctx) {
		ctx.Status(200)
		return
	}

	ctx.Status(400)
}
