package service

// mockgen -source=service/auth-service.go -destination=mock_service/auth-service.go

import (
	"github.com/gin-gonic/gin"
	"github.com/kuritaeiji/todo-gin-back/config"
	"github.com/kuritaeiji/todo-gin-back/dto"
	"github.com/kuritaeiji/todo-gin-back/repository"
)

type AuthService interface {
	Login(*gin.Context) (string, error)
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

// test
func TestNewAuthService(userRepository repository.UserRepository, jwtService JWTService) AuthService {
	return &authService{
		userRepository: userRepository,
		jwtService:     jwtService,
	}
}
