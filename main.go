package main

import (
	"github.com/kuritaeiji/todo-gin-back/config"
	"github.com/kuritaeiji/todo-gin-back/db"
	"github.com/kuritaeiji/todo-gin-back/server"
	"github.com/kuritaeiji/todo-gin-back/validators"
)

func main() {
	config.Init()
	db.Init()
	defer db.CloseDB()
	validators.Init()
	server.Init()
}
