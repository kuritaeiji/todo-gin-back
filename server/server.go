package server

import "github.com/gin-gonic/gin"

func Init() {
	router := gin.Default()
	router.GET("", func(ctx *gin.Context) {
		ctx.String(200, "Hello, World")
	})

	router.Run()
}
