package service_test

import (
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/golang/mock/gomock"
	"github.com/kuritaeiji/todo-gin-back/config"
	"github.com/kuritaeiji/todo-gin-back/factory"
	"github.com/kuritaeiji/todo-gin-back/mock_gateway"
	"github.com/kuritaeiji/todo-gin-back/mock_repository"
	"github.com/kuritaeiji/todo-gin-back/mock_service"
	"github.com/kuritaeiji/todo-gin-back/model"
	"github.com/kuritaeiji/todo-gin-back/service"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type AuthServiceTestSuite struct {
	suite.Suite
	service            service.AuthService
	userRepositoryMock *mock_repository.MockUserRepository
	jwtServiceMock     *mock_service.MockJWTService
	oauthGatewayMock   *mock_gateway.MockOauthGateway
	rec                *httptest.ResponseRecorder
	ctx                *gin.Context
}

func (suite *AuthServiceTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	config.Init()
}

func (suite *AuthServiceTestSuite) SetupTest() {
	suite.userRepositoryMock = mock_repository.NewMockUserRepository(gomock.NewController(suite.T()))
	suite.jwtServiceMock = mock_service.NewMockJWTService(gomock.NewController(suite.T()))
	suite.oauthGatewayMock = mock_gateway.NewMockOauthGateway(gomock.NewController(suite.T()))
	suite.service = service.TestNewAuthService(suite.userRepositoryMock, suite.jwtServiceMock, suite.oauthGatewayMock)
	suite.rec = httptest.NewRecorder()
	suite.ctx, _ = gin.CreateTestContext(suite.rec)
}

func TestAuthService(t *testing.T) {
	suite.Run(t, new(AuthServiceTestSuite))
}

func (suite *AuthServiceTestSuite) TestSuccessLogin() {
	var userConfig factory.UserConfig
	user := factory.NewUser(&userConfig)
	tokenString := factory.CreateAccessToken(user)
	suite.userRepositoryMock.EXPECT().FindByEmail(userConfig.Email).Return(user, nil)
	suite.jwtServiceMock.EXPECT().CreateJWT(user, service.DayFromNowAccessToken).Return(tokenString)

	req := httptest.NewRequest("POST", "/login", factory.CreateUserRequestBody(&userConfig))
	req.Header.Add("Content-Type", binding.MIMEJSON)
	suite.ctx.Request = req
	token, err := suite.service.Login(suite.ctx)

	suite.Equal(tokenString, token)
	suite.Nil(err)
}

func (suite *AuthServiceTestSuite) TestBadLoginWithCannotBind() {
	req := httptest.NewRequest("POST", "/login", nil)
	req.Header.Add("Content-Type", binding.MIMEJSON)
	suite.ctx.Request = req
	_, err := suite.service.Login(suite.ctx)

	suite.Error(err)
}

func (suite *AuthServiceTestSuite) TestBadLoginWithRecordNotFound() {
	var userConfig factory.UserConfig
	req := httptest.NewRequest("POST", "/login", factory.CreateUserRequestBody(&userConfig))
	req.Header.Add("Content-Type", binding.MIMEJSON)
	suite.ctx.Request = req
	suite.userRepositoryMock.EXPECT().FindByEmail(userConfig.Email).Return(model.User{}, gorm.ErrRecordNotFound)
	_, err := suite.service.Login(suite.ctx)

	suite.Equal(gorm.ErrRecordNotFound, err)
}

func (suite *AuthServiceTestSuite) TestBadLoginWithPasswordAuthenticationError() {
	var userConfig factory.UserConfig
	user := factory.NewUser(&userConfig)
	suite.userRepositoryMock.EXPECT().FindByEmail(user.Email).Return(user, nil)

	userConfig.Password = "invalid password"
	req := httptest.NewRequest("POST", "/login", factory.CreateUserRequestBody(&userConfig))
	req.Header.Add("Content-Type", binding.MIMEJSON)
	suite.ctx.Request = req
	_, err := suite.service.Login(suite.ctx)

	suite.Equal(config.PasswordAuthenticationError, err)
}

func (suite *AuthServiceTestSuite) TestSuccessGoogle() {
	const authURL = "https://example.com/auth"
	providerConfig := &oidc.ProviderConfig{AuthURL: authURL}
	provider := providerConfig.NewProvider(suite.ctx)
	suite.oauthGatewayMock.EXPECT().SearchProvider(suite.ctx).Return(provider, nil)

	uri, state, err := suite.service.Google(suite.ctx)
	u, uerr := url.Parse(uri)
	if uerr != nil {
		suite.Fail("url parse error")
	}

	sampleUrl, envUrlError := url.Parse(authURL)
	if envUrlError != nil {
		suite.Fail("env url parse error")
	}

	suite.Equal(sampleUrl.Scheme, u.Scheme)
	suite.Equal(sampleUrl.Hostname(), u.Hostname())
	suite.Equal(sampleUrl.Path, u.Path)
	suite.Equal(os.Getenv("CLIENT_ID"), u.Query()["client_id"][0])
	suite.Equal(os.Getenv("REDIRECT_URL"), u.Query()["redirect_uri"][0])
	suite.Equal("code", u.Query()["response_type"][0])
	suite.Equal("openid", u.Query()["scope"][0])
	suite.Equal(u.Query()["state"][0], state)
	suite.Nil(err)
}

// func (suite *AuthServiceTestSuite) TestSuccessGoogleLogin() {
// 	// stateの検証
// 	const state = "state"
// 	const code = "code"
// 	suite.ctx.SetCookie("state", state, 1000, "", "", false, false)
// 	body := gin.H{
// 		"state": state,
// 		"code":  "code",
// 	}
// 	json, _ := json.Marshal(body)
// 	req := httptest.NewRequest("POST", "/api/google/login", strings.NewReader(string(json)))
// 	suite.ctx.Request = req

// 	// トークンエンドポイントにリクエスト
// 	providerConfig := &oidc.ProviderConfig{}
// 	provider := providerConfig.NewProvider(suite.ctx)
// 	oauth2Config := service.CreateOauth2Config(provider)
// 	const tokenID = "fjoajfhofa1jojo"
// 	reflect.ValueOf(oauth2).MethodByName()
// 	// internalToken := internal.Token{Raw: map[string]interface{"token_id": tokenID}}
// 	suite.oauthGatewayMock.EXPECT().RequestTokenEndpoint(oauth2Config, suite.ctx, code)

// }
