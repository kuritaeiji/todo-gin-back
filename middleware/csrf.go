package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/kuritaeiji/todo-gin-back/config"
)

type CsrfMiddleware interface {
	ConfirmRequestHeader(*gin.Context)
}

type csrfMiddleware struct{}

func NewCsrfMiddleware() CsrfMiddleware {
	return &csrfMiddleware{}
}

// 必ずプリフライトリクエストをさせるためにカスタムヘッダーが付与されていないリクエストは弾く
func (m *csrfMiddleware) ConfirmRequestHeader(ctx *gin.Context) {
	if ctx.GetHeader(config.CsrfCustomHeader["key"]) != config.CsrfCustomHeader["value"] {
		ctx.AbortWithStatus(403)
		return
	}

	ctx.Next()
}
