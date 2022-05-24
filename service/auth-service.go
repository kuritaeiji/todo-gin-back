package service

// mockgen -source=service/auth-service.go -destination=mock_service/auth-service.go

import (
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/kuritaeiji/todo-gin-back/config"
	"github.com/kuritaeiji/todo-gin-back/dto"
	"github.com/kuritaeiji/todo-gin-back/gateway"
	"github.com/kuritaeiji/todo-gin-back/repository"
	"golang.org/x/oauth2"
)

type AuthService interface {
	Login(*gin.Context) (string, error)
	Google(*gin.Context) (string, string, error)
	GoogleLogin(*gin.Context) (string, error)
}

type authService struct {
	dto            dto.Auth
	dtoOauth       dto.Oauth
	userRepository repository.UserRepository
	jwtService     JWTService
	oauthGateway   gateway.OauthGateway
}

func NewAuthService() AuthService {
	return &authService{
		dto:            dto.Auth{},
		userRepository: repository.NewUserRepository(),
		jwtService:     NewJWTService(),
		oauthGateway:   gateway.NewOauthGateway(),
	}
}

func (s *authService) Login(ctx *gin.Context) (string, error) {
	err := ctx.ShouldBindJSON(&s.dto)
	if err != nil {
		return "", err
	}

	user, err := s.userRepository.FindByEmail(s.dto.Email)
	if err != nil {
		return "", err
	}

	if !user.Authenticate(s.dto.Password) {
		return "", config.PasswordAuthenticationError
	}

	return s.jwtService.CreateJWT(user, DayFromNowAccessToken), nil
}

func (s *authService) Google(ctx *gin.Context) (string, string, error) {
	state := config.MakeRandomStr(20)

	provider, err := s.oauthGateway.SearchProvider(ctx)
	if err != nil {
		return "", "", err
	}

	oauth2Config := CreateOauth2Config(provider)
	return oauth2Config.AuthCodeURL(state), state, nil
}

func (s *authService) GoogleLogin(ctx *gin.Context) (string, error) {
	// stateの検証
	cookieState, err := ctx.Cookie(config.StateCookieKey)
	if err != nil {
		return "", err
	}

	ctx.ShouldBindJSON(&s.dtoOauth)
	if cookieState != s.dtoOauth.State {
		return "", config.CsrfError
	}

	provider, err := s.oauthGateway.SearchProvider(ctx)
	if err != nil {
		return "", err
	}

	// トークンエンドポイントにリクエスト
	oauth2Config := CreateOauth2Config(provider)
	oauth2Token, err := s.oauthGateway.RequestTokenEndpoint(oauth2Config, ctx, s.dtoOauth.Code)
	if err != nil {
		return "", err
	}

	// id_tokenの取り出し
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return "", config.StandardError
	}

	// id_tokenの検証
	idToken, err := s.oauthGateway.VerifyIDToken(ctx, provider, rawIDToken)
	if err != nil {
		return "", config.StandardError
	}

	// open_idによるユーザーの作成もしくは探索
	user, err := s.userRepository.FindOrCreateByOpenID(idToken.Subject)
	if err != nil {
		return "", err
	}

	// jwtを作成
	return s.jwtService.CreateJWT(user, DayFromNowAccessToken), nil
}

func CreateOauth2Config(provider *oidc.Provider) oauth2.Config {
	return oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  os.Getenv("REDIRECT_URL"),
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID},
	}
}

// test
func TestNewAuthService(userRepository repository.UserRepository, jwtService JWTService, oauthGateway gateway.OauthGateway) AuthService {
	return &authService{
		userRepository: userRepository,
		jwtService:     jwtService,
		oauthGateway:   oauthGateway,
	}
}
