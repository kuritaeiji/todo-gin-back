package controller

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kuritaeiji/todo-gin-back/config"
	"github.com/kuritaeiji/todo-gin-back/service"
	"gorm.io/gorm"
)

type AuthController interface {
	Login(*gin.Context)       // GET /api/login
	Google(*gin.Context)      // GET /api/google
	GoogleLogin(*gin.Context) // POST /api/google/login
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
		ctx.AbortWithStatus(500)
		return
	}

	ctx.JSON(200, gin.H{
		"token": tokenString,
	})
}

func (c *authController) Google(ctx *gin.Context) {
	url, state, err := c.service.Google(ctx)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}

	var secure bool
	if gin.Mode() == gin.ReleaseMode {
		secure = true
	}

	ctx.SetCookie(config.StateCookieKey, state, 3600, "", os.Getenv("DOMAIN"), secure, true)

	ctx.JSON(200, gin.H{
		"url":   url,
		"state": state,
	})
}

func (c *authController) GoogleLogin(ctx *gin.Context) {
	tokenString, err := c.service.GoogleLogin(ctx)
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(500)
	}

	ctx.JSON(200, gin.H{
		"token": tokenString,
	})
}

// test
func TestNewAuthController(service service.AuthService) AuthController {
	return &authController{
		service: service,
	}
}
