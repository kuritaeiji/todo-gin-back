package request_test

import (
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
	"github.com/kuritaeiji/todo-gin-back/mock_service"
	"github.com/kuritaeiji/todo-gin-back/model"
	"github.com/kuritaeiji/todo-gin-back/server"
	"github.com/kuritaeiji/todo-gin-back/validators"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/stretchr/testify/suite"
)

type UserRequestTestSuite struct {
	suite.Suite
	router *gin.Engine
	mock   *mock_service.MockEmailClient
	rec    *httptest.ResponseRecorder
}

func (suite *UserRequestTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	config.Init()
	validators.Init()
	db.Init()
}

func (suite *UserRequestTestSuite) SetupTest() {
	emailClientMock := mock_service.NewMockEmailClient(gomock.NewController(suite.T()))
	suite.router = server.TestRouterSetup(emailClientMock)
	suite.mock = emailClientMock
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
	doFunc := func(msg *mail.SGMailV3) {
		suite.Equal(os.Getenv("FROM_EMAIL_NAME"), msg.From.Name)
		suite.Equal(os.Getenv("FROM_EMAIL_ADDRESS"), msg.From.Address)
		suite.Equal("アカウント有効化リンク", msg.Subject)
		suite.Equal(email, msg.Personalizations[0].To[0].Address)
		suite.Contains(msg.Content[0].Value, fmt.Sprintf(`<a href="%v/activate?token=`, os.Getenv("FRONT_ORIGIN")))
		re, _ := regexp.Compile(`token=(.+)">`)
		tokenString = re.FindStringSubmatch(msg.Content[0].Value)[0]
	}
	suite.mock.EXPECT().Send(gomock.Any()).Return(&rest.Response{}, nil).Do(doFunc)

	// TODO tokenが正しいかテストする
	println(tokenString)

	suite.router.ServeHTTP(suite.rec, req)

	suite.Equal(200, suite.rec.Code)
}

func (suite *UserRequestTestSuite) TestBadCreateWithInvalid() {
	bodyReader := strings.NewReader(`{"email":"","password":""}`)
	req := httptest.NewRequest("POST", "/users", bodyReader)
	req.Header.Add("Content-Type", binding.MIMEJSON)

	suite.router.ServeHTTP(suite.rec, req)
	suite.Equal(400, suite.rec.Code)
}

func (suite *UserRequestTestSuite) TestBadCreateWithNotUnique() {
	email := "user@example.com"
	password := "Password1010"
	db.GetDB().Create(&model.User{Email: email, PasswordDigest: "pass"})

	bodyReader := strings.NewReader(fmt.Sprintf(`{"email":"%v","password":"%v"}`, email, password))
	req := httptest.NewRequest("POST", "/users", bodyReader)
	req.Header.Add("Content-Type", binding.MIMEJSON)

	suite.router.ServeHTTP(suite.rec, req)
	suite.Equal(400, suite.rec.Code)
}

func (suite *UserRequestTestSuite) TestSuccessUnique() {
	req := httptest.NewRequest("GET", "/users/unique?email=email", nil)
	suite.router.ServeHTTP(suite.rec, req)
	suite.Equal(200, suite.rec.Code)
}

func (suite *UserRequestTestSuite) TestBadUnique() {
	email := "user@example.com"
	db.GetDB().Create(&model.User{Email: email, PasswordDigest: "pass"})
	req := httptest.NewRequest("GET", fmt.Sprintf("/users/unique?email=%v", email), nil)
	suite.router.ServeHTTP(suite.rec, req)
	suite.Equal(400, suite.rec.Code)
}
