package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/kuritaeiji/todo-gin-back/config"
	"github.com/kuritaeiji/todo-gin-back/model"
	"github.com/kuritaeiji/todo-gin-back/service"
	"gorm.io/gorm"
)

type listController struct {
	service service.ListService
}

type ListController interface {
	Index(*gin.Context)   // GET /api/lists
	Create(*gin.Context)  // POST /api/lists
	Update(*gin.Context)  // PUT /api/lists/:id
	Destroy(*gin.Context) // DELETE /api/lists/:id
}

func NewListController() ListController {
	return &listController{service: service.NewListService()}
}

func (c *listController) Index(ctx *gin.Context) {
	lists, err := c.service.Index(ctx)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}

	ctx.JSON(200, model.ToJsonListSlice(lists))
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

	ctx.JSON(200, list.ToJson())
}

func (c *listController) Update(ctx *gin.Context) {
	list, err := c.service.Update(ctx)

	if err == gorm.ErrRecordNotFound {
		ctx.AbortWithStatusJSON(config.RecordNotFoundErrorResponse.Code, config.RecordNotFoundErrorResponse.Json)
		return
	}

	if err == config.ForbiddenError {
		ctx.AbortWithStatusJSON(config.ForbiddenErrorResponse.Code, config.ForbiddenErrorResponse.Json)
		return
	}

	if _, ok := err.(validator.ValidationErrors); ok {
		ctx.AbortWithStatusJSON(config.ValidationErrorResponse.Code, config.ValidationErrorResponse.Json)
		return
	}

	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}

	ctx.JSON(200, list.ToJson())
}

func (c *listController) Destroy(ctx *gin.Context) {
	err := c.service.Destroy(ctx)

	if err == gorm.ErrRecordNotFound {
		ctx.AbortWithStatusJSON(config.RecordNotFoundErrorResponse.Code, config.RecordNotFoundErrorResponse.Json)
		return
	}

	if err == config.ForbiddenError {
		ctx.AbortWithStatusJSON(config.ForbiddenErrorResponse.Code, config.ForbiddenErrorResponse.Json)
		return
	}

	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}

	ctx.Status(200)
}

// test用
func TestNewListController(listService service.ListService) ListController {
	return &listController{service: listService}
}
