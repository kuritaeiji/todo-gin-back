package request_test

import (
	"encoding/json"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/kuritaeiji/todo-gin-back/config"
	"github.com/kuritaeiji/todo-gin-back/controller"
	"github.com/kuritaeiji/todo-gin-back/db"
	"github.com/kuritaeiji/todo-gin-back/dto"
	"github.com/kuritaeiji/todo-gin-back/factory"
	"github.com/kuritaeiji/todo-gin-back/model"
	"github.com/kuritaeiji/todo-gin-back/server"
	"github.com/kuritaeiji/todo-gin-back/service"
	"github.com/kuritaeiji/todo-gin-back/validators"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type AuthRequestTestSuite struct {
	suite.Suite
	router *gin.Engine
	rec    *httptest.ResponseRecorder
	db     *gorm.DB
}

func (suite *AuthRequestTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	config.Init()
	db.Init()
	validators.Init()
	suite.router = server.RouterSetup(controller.NewUserController())
	suite.db = db.GetDB()
}

func (suite *AuthRequestTestSuite) SetupTest() {
	suite.rec = httptest.NewRecorder()
}

func (suite *AuthRequestTestSuite) TearDownSuite() {
	db.CloseDB()
}

func (suite *AuthRequestTestSuite) TearDownTest() {
	db.DeleteAll()
}

func TestAuthRequest(t *testing.T) {
	suite.Run(t, new(AuthRequestTestSuite))
}

var (
	email    = "user@example.com"
	password = "Password1010"
)

func (suite *AuthRequestTestSuite) TestSuccessLogin() {
	dtoUser := dto.User{Email: email, Password: password}
	var user model.User
	dtoUser.Transfer(&user)
	suite.db.Create(&user)

	body := map[string]string{
		"email":    email,
		"password": password,
	}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/login", strings.NewReader(string(bodyBytes)))
	req.Header.Add("Content-Type", binding.MIMEJSON)
	suite.router.ServeHTTP(suite.rec, req)

	suite.Equal(200, suite.rec.Code)
	re, _ := regexp.Compile(`"token":"(.+)"`)
	tokenString := re.FindStringSubmatch(suite.rec.Body.String())[1]
	claim, err := service.NewJWTService().VerifyJWT(tokenString)
	suite.Equal(user.ID, claim.ID)
	suite.InEpsilon(time.Now().AddDate(0, 0, service.DayFromNowAccessToken).Unix(), claim.ExpiresAt, 30)
	suite.Nil(err)
}

func (suite *AuthRequestTestSuite) TestBadLoginWithAreadyLogin() {
	user := factory.CreateUser(&factory.UserConfig{Email: email, Password: password})
	tokenString := factory.CreateAccessToken(user)
	req := httptest.NewRequest("POST", "/api/login", nil)
	req.Header.Add("Authorization", "Bearer "+tokenString)
	suite.router.ServeHTTP(suite.rec, req)

	suite.Equal(config.GuestErrorResponse.Code, suite.rec.Code)
	suite.Contains(suite.rec.Body.String(), config.GuestErrorResponse.Json["content"])
}

func (suite *AuthRequestTestSuite) TestBadLoginWithRecordNotFound() {
	body := map[string]string{
		"email": "",
	}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/login", strings.NewReader(string(bodyBytes)))
	req.Header.Add("Content-Type", binding.MIMEJSON)
	suite.router.ServeHTTP(suite.rec, req)

	suite.Equal(config.RecordNotFoundErrorResponse.Code, suite.rec.Code)
	suite.Contains(suite.rec.Body.String(), config.RecordNotFoundErrorResponse.Json["content"])
}

func (suite *AuthRequestTestSuite) TestBadLoginWithPasswordAuthenticationError() {
	dtoUser := dto.User{Email: email, Password: password}
	var user model.User
	dtoUser.Transfer(&user)
	suite.db.Create(&user)

	body := map[string]string{
		"email":    email,
		"password": "pass",
	}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/login", strings.NewReader(string(bodyBytes)))
	req.Header.Add("Content-Type", binding.MIMEJSON)
	suite.router.ServeHTTP(suite.rec, req)

	suite.Equal(config.PasswordAuthenticationErrorResponse.Code, suite.rec.Code)
	suite.Contains(suite.rec.Body.String(), config.PasswordAuthenticationErrorResponse.Json["content"])
}
