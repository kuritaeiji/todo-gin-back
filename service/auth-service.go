package service

// mockgen -source=service/auth-service.go -destination=mock_service/auth-service.go

import (
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
	Google(*gin.Context) (string, error)
}

type authService struct {
	dto            dto.Auth
	userRepository repository.UserRepository
	jwtService     JWTService
}

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

func (s *authService) Google(ctx *gin.Context) (string, error) {
	provider, err := oidc.NewProvider(ctx.Request.Context(), os.Getenv("GOOGLE_OAUTH_URL"))
	if err != nil {
		return "", err
	}

	oauth2Config := oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  os.Getenv("REDIRECT_URL"),
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID},
	}

	return oauth2Config.AuthCodeURL(config.MakeRandomStr(20)), nil
}

// test
func TestNewAuthService(userRepository repository.UserRepository, jwtService JWTService) AuthService {
	return &authService{
		userRepository: userRepository,
		jwtService:     jwtService,
	}
}
