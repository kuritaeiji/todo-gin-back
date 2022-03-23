package server

import (
	"github.com/gin-gonic/gin"
	"github.com/kuritaeiji/todo-gin-back/controller"
	"github.com/kuritaeiji/todo-gin-back/mock_gateway"
	"github.com/kuritaeiji/todo-gin-back/service"
)

func Init() {
	router := routerSetup(controller.NewUserController())
	router.Run()
}

func routerSetup(userController controller.UserController) *gin.Engine {
	r := gin.Default()

	user := r.Group("/users")
	{
		user.POST("", userController.Create)
		user.GET("/unique", userController.IsUnique)
		user.PUT("/activate", userController.Activate)
	}

	return r
}

// test用 sendgridのmailclientをモック化
func TestRouterSetup(emailClientMock *mock_gateway.MockEmailGateway) *gin.Engine {
	con := controller.TestNewUserController(service.NewUserService(), service.TestNewEmailService(
		emailClientMock, service.NewJWTService(),
	))
	return routerSetup(con)
}
