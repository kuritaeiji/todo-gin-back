package factory

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/kuritaeiji/todo-gin-back/db"
	"github.com/kuritaeiji/todo-gin-back/dto"
	"github.com/kuritaeiji/todo-gin-back/model"
)

var defaultTitle = "list title"

type ListConfig struct {
	Title              string
	Index              int
	NotUseDefaultValue bool
}

func (config *ListConfig) setDefaultValue() {
	if config.NotUseDefaultValue {
		return
	}

	if config.Title == "" {
		config.Title = defaultTitle
	}
}

func NewDtoList(config *ListConfig) dto.List {
	config.setDefaultValue()
	return dto.List{Title: config.Title, Index: config.Index}
}

func NewList(config *ListConfig) model.List {
	dtoList := NewDtoList(config)
	var list model.List
	dtoList.Transfer(&list)
	return list
}

func CreateList(config *ListConfig, user model.User) model.List {
	list := NewList(config)
	list.UserID = user.ID
	db.GetDB().Create(&list)
	return list
}

func CreateListRequestBody(config *ListConfig) io.Reader {
	config.setDefaultValue()
	body := map[string]interface{}{
		"title": config.Title,
		"index": config.Index,
	}
	bodyBytes, _ := json.Marshal(body)
	return strings.NewReader(string(bodyBytes))
}
