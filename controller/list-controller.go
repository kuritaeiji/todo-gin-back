package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/kuritaeiji/todo-gin-back/config"
	"github.com/kuritaeiji/todo-gin-back/service"
)

type listController struct {
	service service.ListService
}

type ListController interface {
	Create(*gin.Context)
}

func NewListController() ListController {
	return &listController{service: service.NewListService()}
}

func (c *listController) Create(ctx *gin.Context) {
	list, err := c.service.Create(ctx)

	if _, ok := err.(validator.ValidationErrors); ok {
		ctx.AbortWithStatusJSON(config.ValidationErrorResponse.Code, config.ValidationErrorResponse.Json)
		return
	}

	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}

	ctx.JSON(200, list)
}

// test用
func TestNewListController(listService service.ListService) ListController {
	return &listController{service: listService}
}
