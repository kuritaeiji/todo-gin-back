package middleware

import (
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kuritaeiji/todo-gin-back/config"
)

func NewCorsMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: []string{os.Getenv("FRONT_ORIGIN"), "https://todo-gin.ml"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Access-Control-Allow-Headers",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"Authorization",
			config.CsrfCustomHeader["key"],
		},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour,
	})
}
