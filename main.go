package main

import (
	"github.com/kuritaeiji/todo-gin-back/config"
	"github.com/kuritaeiji/todo-gin-back/db"
	"github.com/kuritaeiji/todo-gin-back/server"
)

func main() {
	config.Init()
	db.Init()
	defer db.CloseDB()
	server.Init()
}
