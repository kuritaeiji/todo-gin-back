package controller_test

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/kuritaeiji/todo-gin-back/db"
	"github.com/kuritaeiji/todo-gin-back/model"
	"github.com/kuritaeiji/todo-gin-back/server"
	"github.com/kuritaeiji/todo-gin-back/validators"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var (
	router   *gin.Engine
	database *gorm.DB
)

func TestMain(m *testing.M) {
	validators.Init()
	db.TestInit()
	database = db.GetDB()
	defer db.CloseDB()
	router = server.RouterSetUp()

	m.Run()
}

func TestCreate(t *testing.T) {
	assert := assert.New(t)

	tests := map[string]struct {
		body       string
		code       int
		userLength int
	}{
		"normal":  {body: `{"email":"user@example.com","password":"Password1010"}`, code: 200, userLength: 1},
		"invalid": {body: `{"email":"","password":""}`, code: 400, userLength: 0},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			bodyReader := strings.NewReader(test.body)
			req := httptest.NewRequest("POST", "/users", bodyReader)
			req.Header.Add("Content-Type", binding.MIMEJSON)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			var count int64
			database.Model(&model.User{}).Count(&count)

			assert.Equal(test.code, rec.Code)
			assert.Equal(test.userLength, int(count))

			if testName == "normal" {
				var user model.User
				json.Unmarshal(rec.Body.Bytes(), &user)
				assert.Equal("user@example.com", user.Email)
			}
			db.DeleteAll()
		})
	}
}

func TestIsUniqueEmail(t *testing.T) {
	assert := assert.New(t)

	email := "user@example.com"
	tests := map[string]struct {
		callback func()
		code     int
	}{
		"unique": {func() {}, 200},
		"not unique": {func() {
			database.Create(&model.User{Email: email, PasswordDigest: "pass"})
		}, 400},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			test.callback()
			req := httptest.NewRequest("GET", fmt.Sprintf("/users/unique-email?email=%v", email), nil)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			assert.Equal(rec.Code, test.code)
			db.DeleteAll()
		})
	}
}
