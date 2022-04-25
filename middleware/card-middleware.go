package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/kuritaeiji/todo-gin-back/config"
	"github.com/kuritaeiji/todo-gin-back/service"
	"gorm.io/gorm"
)

type cardMiddleware struct {
	service service.CardMiddlewareService
}

type CardMiddleware interface {
	Authorize(*gin.Context)
}

func NewCardMiddleware() CardMiddleware {
	return &cardMiddleware{service: service.NewCardMiddlewareService()}
}

func (m *cardMiddleware) Authorize(ctx *gin.Context) {
	card, err := m.service.Authorize(ctx)
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

	ctx.Set(config.CardKey, card)
	ctx.Next()
}

// test
func TestNewCardMiddleware(cardMiddlewareService service.CardMiddlewareService) CardMiddleware {
	return &cardMiddleware{service: cardMiddlewareService}
}
