package request_test

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/golang/mock/gomock"
	"github.com/kuritaeiji/todo-gin-back/config"
	"github.com/kuritaeiji/todo-gin-back/db"
	"github.com/kuritaeiji/todo-gin-back/mock_gateway"
	"github.com/kuritaeiji/todo-gin-back/model"
	"github.com/kuritaeiji/todo-gin-back/server"
	"github.com/kuritaeiji/todo-gin-back/service"
	"github.com/kuritaeiji/todo-gin-back/validators"
	"github.com/stretchr/testify/suite"
)

type UserRequestTestSuite struct {
	suite.Suite
	router *gin.Engine
	mock   *mock_gateway.MockEmailGateway
	rec    *httptest.ResponseRecorder
}

func (suite *UserRequestTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	config.Init()
	validators.Init()
	db.Init()
}

func (suite *UserRequestTestSuite) SetupTest() {
	emailGatewayMock := mock_gateway.NewMockEmailGateway(gomock.NewController(suite.T()))
	suite.router = server.TestRouterSetup(emailGatewayMock)
	suite.mock = emailGatewayMock
	suite.rec = httptest.NewRecorder()
}

func (suite *UserRequestTestSuite) TearDownSuite() {
	db.CloseDB()
}

func (suite *UserRequestTestSuite) TearDownTest() {
	db.DeleteAll()
}

func TestUserRequestSuite(t *testing.T) {
	suite.Run(t, &UserRequestTestSuite{})
}

func (suite *UserRequestTestSuite) TestSuccessCreate() {
	email := "user@example.com"
	password := "Password1010"
	bodyReader := strings.NewReader(fmt.Sprintf(`{"email":"%v","password":"%v"}`, email, password))
	req := httptest.NewRequest("POST", "/users", bodyReader)
	req.Header.Add("Content-Type", binding.MIMEJSON)

	var tokenString string
	doFunc := func(email, subject, htmlString string) {
		suite.Contains(htmlString, fmt.Sprintf(`<a href="%v/activate?token=`, os.Getenv("FRONT_ORIGIN")))
		re, _ := regexp.Compile(`token=(.+)">`)
		tokenString = re.FindStringSubmatch(htmlString)[1]
		claim, err := service.NewJWTService().VerifyJWT(tokenString)
		suite.Nil(err)
		var user model.User
		db.GetDB().First(&user)
		suite.Equal(user.ID, claim.ID)
	}
	suite.mock.EXPECT().Send(email, "アカウント有効化リンク", gomock.Any()).Return(nil).Do(doFunc)

	suite.router.ServeHTTP(suite.rec, req)
	suite.Equal(200, suite.rec.Code)
}

func (suite *UserRequestTestSuite) TestBadCreateWithNotGuest() {
	req := httptest.NewRequest("POST", "/users", nil)
	req.Header.Add("Authorization", "Bearer token")
	suite.router.ServeHTTP(suite.rec, req)

	suite.Equal(config.GuestErrorResponse.Code, suite.rec.Code)
	suite.Contains(suite.rec.Body.String(), config.GuestErrorResponse.Json["content"])
}

func (suite *UserRequestTestSuite) TestBadCreateWithInvalid() {
	bodyReader := strings.NewReader(`{"email":"","password":""}`)
	req := httptest.NewRequest("POST", "/users", bodyReader)
	req.Header.Add("Content-Type", binding.MIMEJSON)

	suite.router.ServeHTTP(suite.rec, req)
	suite.Equal(config.ValidationErrorResponse.Code, suite.rec.Code)
	body := suite.rec.Body.String()
	suite.Contains(body, config.ValidationErrorResponse.Json["content"])
}

func (suite *UserRequestTestSuite) TestBadCreateWithNotUnique() {
	email := "user@example.com"
	password := "Password1010"
	db.GetDB().Create(&model.User{Email: email, PasswordDigest: "pass"})

	bodyReader := strings.NewReader(fmt.Sprintf(`{"email":"%v","password":"%v"}`, email, password))
	req := httptest.NewRequest("POST", "/users", bodyReader)
	req.Header.Add("Content-Type", binding.MIMEJSON)

	suite.router.ServeHTTP(suite.rec, req)
	suite.Equal(config.UniqueUserErrorResponse.Code, suite.rec.Code)
}

func (suite *UserRequestTestSuite) TestBadCreateWithEmailClientError() {
	email := "user@example.com"
	body := map[string]string{
		"email":    email,
		"password": "Password1010",
	}
	bodyJson, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	req := httptest.NewRequest("POST", "/users", strings.NewReader(string(bodyJson)))
	req.Header.Add("Content-Type", binding.MIMEJSON)
	suite.mock.EXPECT().Send(email, "アカウント有効化リンク", gomock.Any()).Return(config.EmailClientError)
	suite.router.ServeHTTP(suite.rec, req)

	suite.Equal(config.EmailClientErrorResponse.Code, suite.rec.Code)
	suite.Contains(suite.rec.Body.String(), config.EmailClientErrorResponse.Json["content"])
}

func (suite *UserRequestTestSuite) TestSuccessUnique() {
	req := httptest.NewRequest("GET", "/users/unique?email=email", nil)
	suite.router.ServeHTTP(suite.rec, req)
	suite.Equal(200, suite.rec.Code)
}

func (suite *UserRequestTestSuite) TestBadIsUniqueWithNotGuest() {
	req := httptest.NewRequest("GET", "/users/unique", nil)
	req.Header.Add("Authorization", "Bearer token")
	suite.router.ServeHTTP(suite.rec, req)

	suite.Equal(config.GuestErrorResponse.Code, suite.rec.Code)
	suite.Contains(suite.rec.Body.String(), config.GuestErrorResponse.Json["content"])
}

func (suite *UserRequestTestSuite) TestBadUnique() {
	email := "user@example.com"
	db.GetDB().Create(&model.User{Email: email, PasswordDigest: "pass"})
	req := httptest.NewRequest("GET", fmt.Sprintf("/users/unique?email=%v", email), nil)
	suite.router.ServeHTTP(suite.rec, req)
	suite.Equal(400, suite.rec.Code)
}

func (suite *UserRequestTestSuite) TestSuccessActivate() {
	id := 1
	user := model.User{ID: id, Email: "mail", PasswordDigest: "pass"}
	db.GetDB().Create(&user)
	tokenString := service.NewJWTService().CreateJWT(user, 1)

	req := httptest.NewRequest("PUT", fmt.Sprintf("/users/activate?token=%s", tokenString), nil)
	suite.router.ServeHTTP(suite.rec, req)

	suite.Equal(200, suite.rec.Code)
	db.GetDB().First(&user)
	suite.True(user.Activated)
}

func (suite *UserRequestTestSuite) TestBadActivateWithNotGuest() {
	req := httptest.NewRequest("PUT", "/users/activate", nil)
	req.Header.Add("Authorization", "Bearer token")
	suite.router.ServeHTTP(suite.rec, req)

	suite.Equal(config.GuestErrorResponse.Code, suite.rec.Code)
	suite.Contains(suite.rec.Body.String(), config.GuestErrorResponse.Json["content"])
}

func (suite *UserRequestTestSuite) TestBadActivateWithExpiredJWT() {
	tokenString := service.NewJWTService().CreateJWT(model.User{}, -1)
	req := httptest.NewRequest("PUT", fmt.Sprintf("/users/activate?token=%s", tokenString), nil)
	suite.router.ServeHTTP(suite.rec, req)

	suite.Equal(config.JWTExpiredErrorResponse.Code, suite.rec.Code)
	body := suite.rec.Body.String()
	suite.Contains(body, config.JWTExpiredErrorResponse.Json["content"])
}

func (suite *UserRequestTestSuite) TestBadActivateWithInvalidJWT() {
	tokenString := service.NewJWTService().CreateJWT(model.User{}, 1) + "invalid"
	req := httptest.NewRequest("PUT", fmt.Sprintf("/users/activate?token=%s", tokenString), nil)
	suite.router.ServeHTTP(suite.rec, req)

	suite.Equal(config.JWTValidationErrorResponse.Code, suite.rec.Code)
	body := suite.rec.Body.String()
	suite.Contains(body, config.JWTValidationErrorResponse.Json["content"])
}

func (suite *UserRequestTestSuite) TestBadActivateWithRecordNotFound() {
	tokenString := service.NewJWTService().CreateJWT(model.User{}, 1)
	req := httptest.NewRequest("PUT", fmt.Sprintf("/users/activate?token=%s", tokenString), nil)
	suite.router.ServeHTTP(suite.rec, req)

	suite.Equal(config.RecordNotFoundErrorResponse.Code, suite.rec.Code)
	body := suite.rec.Body.String()
	suite.Contains(body, config.RecordNotFoundErrorResponse.Json["content"])
}

func (suite *UserRequestTestSuite) TestBadActivateWithAlreadyActivatedUser() {
	user := model.User{Email: "email", PasswordDigest: "pass", ID: 1, Activated: true}
	db.GetDB().Create(&user)
	tokenString := service.NewJWTService().CreateJWT(user, service.DayFromNowActivateUserToken)
	req := httptest.NewRequest("PUT", fmt.Sprintf("/users/activate?token=%s", tokenString), nil)
	suite.router.ServeHTTP(suite.rec, req)

	suite.Equal(config.AlreadyActivatedUserErrorResponse.Code, suite.rec.Code)
	body := suite.rec.Body.String()
	suite.Contains(body, config.AlreadyActivatedUserErrorResponse.Json["content"])
}
