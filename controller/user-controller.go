package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/kuritaeiji/todo-gin-back/config"
	"github.com/kuritaeiji/todo-gin-back/service"
	"gorm.io/gorm"
)

type UserController interface {
	Create(ctx *gin.Context)
	IsUnique(ctx *gin.Context)
	Activate(ctx *gin.Context)
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
	if _, ok := err.(validator.ValidationErrors); ok {
		ctx.JSON(config.ValidationErrorResponse.Code, config.ValidationErrorResponse.Json)
		return
	}

	if err == config.UniqueUserError {
		ctx.JSON(config.UniqueUserErrorResponse.Code, config.UniqueUserErrorResponse.Json)
		return
	}

	if err := c.emailService.ActivationUserEmail(user); err != nil {
		ctx.JSON(config.EmailClientErrorResponse.Code, config.EmailClientErrorResponse.Json)
		return
	}

	if err != nil {
		ctx.AbortWithStatus(500)
		gin.DefaultWriter.Write([]byte(err.Error()))
		return
	}

	ctx.Status(200)
}

func (c *userController) IsUnique(ctx *gin.Context) {
	result, err := c.service.IsUnique(ctx)
	if !result {
		ctx.JSON(config.UniqueUserErrorResponse.Code, config.UniqueUserErrorResponse.Json)
		return
	}

	if err != nil {
		ctx.AbortWithStatus(500)
		gin.DefaultWriter.Write([]byte(err.Error()))
		return
	}

	ctx.Status(200)
}

func (c *userController) Activate(ctx *gin.Context) {
	err := c.service.Activate(ctx)
	jwtErr, ok := err.(*jwt.ValidationError)
	if ok && jwtErr.Errors == jwt.ValidationErrorExpired {
		ctx.JSON(config.JWTExpiredErrorResponse.Code, config.JWTExpiredErrorResponse.Json)
		return
	}
	if ok {
		ctx.JSON(config.JWTValidationErrorResponse.Code, config.JWTValidationErrorResponse.Json)
		return
	}

	if err == gorm.ErrRecordNotFound {
		ctx.JSON(config.RecordNotFoundErrorResponse.Code, config.RecordNotFoundErrorResponse.Json)
		return
	}

	if err == config.AlreadyActivatedUserError {
		ctx.JSON(config.AlreadyActivatedUserErrorResponse.Code, config.AlreadyActivatedUserErrorResponse.Json)
		return
	}

	if err != nil {
		ctx.AbortWithStatus(500)
		gin.DefaultWriter.Write([]byte(err.Error()))
		return
	}

	ctx.Status(200)
}

// test用
func TestNewUserController(us service.UserService, es service.EmailService) UserController {
	return &userController{us, es}
}
