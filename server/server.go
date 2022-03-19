package server

import (
	"github.com/gin-gonic/gin"
	"github.com/kuritaeiji/todo-gin-back/controller"
)

func Init() {
	router := RouterSetUp()
	router.Run()
}

func RouterSetUp() *gin.Engine {
	r := gin.Default()

	user := r.Group("/users")
	{
		con := controller.NewUserController()
		user.POST("", con.Create)
		user.GET("/unique-email", con.IsUniqueEmail)
	}

	return r
}
