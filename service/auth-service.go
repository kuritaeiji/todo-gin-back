package service

// mockgen -source=service/auth-service.go -destination=mock_service/auth-service.go

import (
	"fmt"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/kuritaeiji/todo-gin-back/config"
	"github.com/kuritaeiji/todo-gin-back/dto"
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
}

var idToken *oidc.IDToken

func NewAuthService() AuthService {
	return &authService{
		dto:            dto.Auth{},
		userRepository: repository.NewUserRepository(),
		jwtService:     NewJWTService(),
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
	provider, err := searchProvider(ctx)
	if err != nil {
		return "", "", err
	}

	oauth2Config := createOauth2Config(provider)
	state := config.MakeRandomStr(20)

	return oauth2Config.AuthCodeURL(state), state, nil
}

func (s *authService) GoogleLogin(ctx *gin.Context) (string, error) {
	cookieState, err := ctx.Cookie("state")
	if err != nil {
		return "", config.CsrfError
	}

	ctx.ShouldBindJSON(&s.dtoOauth)
	if cookieState != s.dtoOauth.State {
		return "", config.CsrfError
	}

	provider, err := searchProvider(ctx)
	if err != nil {
		return "", err
	}

	oauth2Config := createOauth2Config(provider)
	oauth2Token, err := oauth2Config.Exchange(ctx.Request.Context(), s.dtoOauth.Code)
	if rErr, ok := err.(*oauth2.RetrieveError); ok {
		body := string(rErr.Body)
		fmt.Println(body)
	}
	if err != nil {
		return "", err
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return "", config.StandardError
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: os.Getenv("CLIENT_ID")})
	idToken, err = verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return "", config.StandardError
	}

	user, err := s.userRepository.FindOrCreateByOpenID(idToken.Subject)
	if err != nil {
		return "", err
	}

	return s.jwtService.CreateJWT(user, DayFromNowAccessToken), nil
}

func searchProvider(ctx *gin.Context) (*oidc.Provider, error) {
	return oidc.NewProvider(ctx.Request.Context(), os.Getenv("GOOGLE_OAUTH_URL"))
}

func createOauth2Config(provider *oidc.Provider) oauth2.Config {
	return oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  os.Getenv("REDIRECT_URL"),
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID},
	}
}

// test
func TestNewAuthService(userRepository repository.UserRepository, jwtService JWTService) AuthService {
	return &authService{
		userRepository: userRepository,
		jwtService:     jwtService,
	}
}
