package service_test

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kuritaeiji/todo-gin-back/config"
	"github.com/kuritaeiji/todo-gin-back/validators"
	"github.com/stretchr/testify/assert"
)

var (
	assertion *assert.Assertions
	ctx       *gin.Context
	rec       *httptest.ResponseRecorder
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	config.Init()
	validators.Init()
	m.Run()
}
