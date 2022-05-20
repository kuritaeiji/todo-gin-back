package middleware_test

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kuritaeiji/todo-gin-back/config"
	"github.com/kuritaeiji/todo-gin-back/middleware"
	"github.com/stretchr/testify/suite"
)

type CsrfMiddlewareTestSuite struct {
	suite.Suite
	middleware middleware.CsrfMiddleware
	rec        *httptest.ResponseRecorder
	ctx        *gin.Context
}

func (suite *CsrfMiddlewareTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
}

func (suite *CsrfMiddlewareTestSuite) SetupTest() {
	suite.middleware = middleware.NewCsrfMiddleware()
	suite.rec = httptest.NewRecorder()
	suite.ctx, _ = gin.CreateTestContext(suite.rec)
}

func TestCsrfMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(CsrfMiddlewareTestSuite))
}

func (suite *CsrfMiddlewareTestSuite) TestSuccessConfirmRequestHeader() {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add(config.CsrfCustomHeader["key"], config.CsrfCustomHeader["value"])
	suite.ctx.Request = req
	suite.middleware.ConfirmRequestHeader(suite.ctx)

	suite.Equal(200, suite.rec.Code)
}

func (suite *CsrfMiddlewareTestSuite) TestBadConfirmRequestHeader() {
	req := httptest.NewRequest("GET", "/", nil)
	suite.ctx.Request = req
	suite.middleware.ConfirmRequestHeader(suite.ctx)

	suite.Equal(403, suite.rec.Code)
}
