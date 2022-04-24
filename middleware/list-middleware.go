package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/kuritaeiji/todo-gin-back/config"
	"github.com/kuritaeiji/todo-gin-back/service"
	"gorm.io/gorm"
)

type listMiddleware struct {
	service service.ListMiddlewareServive
}

type ListMiddleware interface {
	Authorize(*gin.Context)
}

func NewListMiddleware() ListMiddleware {
	return &listMiddleware{service: service.NewListMiddlewareService()}
}

func (m *listMiddleware) Authorize(ctx *gin.Context) {
	list, err := m.service.Authorize(ctx)
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

	ctx.Set(config.ListKey, list)
	ctx.Next()
}

// test
func TestNewListMiddleware(listMiddlewareService service.ListMiddlewareServive) ListMiddleware {
	return &listMiddleware{service: listMiddlewareService}
}
