package server

import (
	// "bytes"
	"fmt"
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kuritaeiji/todo-gin-back/config"
	"github.com/kuritaeiji/todo-gin-back/controller"
	"github.com/kuritaeiji/todo-gin-back/middleware"
	"github.com/kuritaeiji/todo-gin-back/mock_gateway"
	"github.com/kuritaeiji/todo-gin-back/model"
	"github.com/kuritaeiji/todo-gin-back/service"
)

func Init() {
	router := RouterSetup(controller.NewUserController())
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}

func RouterSetup(userController controller.UserController) *gin.Engine {
	r := router()
	r.Use(middleware.NewCorsMiddleware())
	if gin.Mode() != gin.TestMode {
		r.Use(middleware.NewCsrfMiddleware().ConfirmRequestHeader)
	}

	api := r.Group("/api")

	authMiddleware := middleware.NewAuthMiddleware()
	guest := api.Group("")
	{
		guest.Use(authMiddleware.Guest)

		authCon := controller.NewAuthController()
		guest.POST("/login", authCon.Login)
		guest.GET("/google", authCon.Google)
		guest.POST("/google/login", authCon.GoogleLogin)

		user := guest.Group("/users")
		{
			user.POST("", userController.Create)
			user.GET("/unique", userController.IsUnique)
			user.PUT("/activate", userController.Activate)
		}
	}

	auth := api.Group("")
	{
		auth.Use(authMiddleware.Auth)
		auth.DELETE("/users", userController.Destroy)

		listMiddleware := middleware.NewListMiddleware()
		list := auth.Group("/lists")
		{
			listCon := controller.NewListController()
			list.GET("", listCon.Index)
			list.POST("", listCon.Create)

			listAuth := list.Group("")
			{
				listAuth.Use(listMiddleware.Authorize)
				listAuth.PUT("/:id", listCon.Update)
				listAuth.DELETE("/:id", listCon.Destroy)
				listAuth.PUT("/:id/move", listCon.Move)
			}
		}

		cardCon := controller.NewCardController()
		cardWithListAuth := auth.Group("")
		{
			cardWithListAuth.Use(listMiddleware.Authorize)
			cardWithListAuth.POST("/lists/:listID/cards", cardCon.Create)
		}

		card := auth.Group("/cards")
		{
			cardMiddleware := middleware.NewCardMiddleware()
			card.Use(cardMiddleware.Authorize)
			card.PUT("/:id", cardCon.Update)
			card.DELETE("/:id", cardCon.Destroy)
			card.PUT("/:id/move", cardCon.Move)
		}
	}

	return r
}

func router() *gin.Engine {
	if gin.Mode() == gin.ReleaseMode {
		router := gin.New()
		f, _ := os.Create(fmt.Sprintf("%v/gin.log", os.Getenv("TODO_GIN_WORKDIR")))
		gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
		router.Use(gin.LoggerWithFormatter(func(params gin.LogFormatterParams) string {
			user, _ := params.Keys[config.CurrentUserKey].(model.User)
			return fmt.Sprintf("日時: %s IP: %s リクエスト: [%s] %s ユーザーID: %v レスポンス: %v エラー: %s\n",
				params.TimeStamp.Format("2006-01-02 03:04:05"),
				params.ClientIP,
				params.Method,
				params.Path,
				user.ID,
				params.StatusCode,
				params.ErrorMessage,
			)
		}))
		router.Use(gin.Recovery())
		router.Use(func(ctx *gin.Context) { fmt.Println(ctx.Errors.ByType(gin.ErrorTypePrivate).String()) })
		return router
	}

	return gin.Default()
}

// test用 sendgridのmailclientをモック化
func TestRouterSetup(emailClientMock *mock_gateway.MockEmailGateway) *gin.Engine {
	con := controller.TestNewUserController(service.NewUserService(), service.TestNewEmailService(
		emailClientMock, service.NewJWTService(),
	))
	return RouterSetup(con)
}
