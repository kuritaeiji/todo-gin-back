package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kuritaeiji/todo-gin-back/config"
	"github.com/kuritaeiji/todo-gin-back/db"
	"github.com/kuritaeiji/todo-gin-back/seed"
	"github.com/kuritaeiji/todo-gin-back/server"
	"github.com/kuritaeiji/todo-gin-back/validators"
)

func main() {
	config.Init()
	db.Init()
	defer db.CloseDB()
	validators.Init()
	if gin.Mode() == gin.ReleaseMode {
		seed.CreateSeedData()
	}
	server.Init()
}
