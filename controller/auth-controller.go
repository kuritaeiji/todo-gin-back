package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/kuritaeiji/todo-gin-back/config"
	"github.com/kuritaeiji/todo-gin-back/service"
	"gorm.io/gorm"
)

type AuthController interface {
	Login(*gin.Context)  // GET /api/login
	Google(*gin.Context) // GET /api/google
}

type authController struct {
	service service.AuthService
}

func NewAuthController() AuthController {
	return &authController{
		service: service.NewAuthService(),
	}
}

func (c *authController) Login(ctx *gin.Context) {
	tokenString, err := c.service.Login(ctx)
	if err == gorm.ErrRecordNotFound {
		ctx.JSON(config.RecordNotFoundErrorResponse.Code, config.RecordNotFoundErrorResponse.Json)
		return
	}

	if err == config.PasswordAuthenticationError {
		ctx.JSON(config.PasswordAuthenticationErrorResponse.Code, config.PasswordAuthenticationErrorResponse.Json)
		return
	}

	if err != nil {
		ctx.Status(500)
		return
	}

	ctx.JSON(200, gin.H{
		"token": tokenString,
	})
}

func (c *authController) Google(ctx *gin.Context) {
	url, err := c.service.Google(ctx)
	if err != nil {
		ctx.Status(500)
		return
	}

	ctx.JSON(200, gin.H{
		"url": url,
	})
}

func (c *authController) GoogleCallback(ctx *gin.Context) {

}

// test
func TestNewAuthController(service service.AuthService) AuthController {
	return &authController{
		service: service,
	}
}
