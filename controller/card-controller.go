package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/kuritaeiji/todo-gin-back/config"
	"github.com/kuritaeiji/todo-gin-back/service"
)

type cardController struct {
	service service.CardService
}

type CardController interface {
	Create(*gin.Context)  // POST /api/lists/:listID/cards
	Update(*gin.Context)  // PUT /api/cards/:id
	Destroy(*gin.Context) // DELETE /api/cards/:id
	Move(*gin.Context)    // PUT /api/cards/:id/move
}

func NewCardController() CardController {
	return &cardController{service: service.NewCardService()}
}

func (c *cardController) Create(ctx *gin.Context) {
	card, err := c.service.Create(ctx)

	if _, ok := err.(validator.ValidationErrors); ok {
		ctx.AbortWithStatusJSON(config.ValidationErrorResponse.Code, config.ValidationErrorResponse.Json)
		return
	}

	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}

	ctx.JSON(200, card.ToJson())
}

func (c *cardController) Update(ctx *gin.Context) {
	card, err := c.service.Update(ctx)

	if _, ok := err.(validator.ValidationErrors); ok {
		ctx.AbortWithStatusJSON(config.ValidationErrorResponse.Code, config.ValidationErrorResponse.Json)
		return
	}

	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}

	ctx.JSON(200, card.ToJson())
}

func (c *cardController) Destroy(ctx *gin.Context) {
	err := c.service.Destroy(ctx)

	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}

	ctx.Status(200)
}

func (c *cardController) Move(ctx *gin.Context) {
	err := c.service.Move(ctx)

	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}

	ctx.Status(200)
}

// test
func TestNewCardController(cardService service.CardService) CardController {
	return &cardController{service: cardService}
}
