package service

// mockgen -source=service/user-service.go -destination=./mock_service/user-service.go

import (
	"github.com/gin-gonic/gin"
	"github.com/kuritaeiji/todo-gin-back/config"
	"github.com/kuritaeiji/todo-gin-back/dto"
	"github.com/kuritaeiji/todo-gin-back/model"
	"github.com/kuritaeiji/todo-gin-back/repository"
)

type UserService interface {
	Create(*gin.Context) (model.User, error)
	IsUnique(*gin.Context) (bool, error)
	Activate(*gin.Context) error
	Destroy(*gin.Context) error
}

type userService struct {
	jwtService JWTService
	repository repository.UserRepository
	dto        dto.User
}

func NewUserService() UserService {
	return &userService{
		jwtService: NewJWTService(),
		repository: repository.NewUserRepository(),
		dto:        dto.User{},
	}
}

func (s *userService) Create(ctx *gin.Context) (model.User, error) {
	if err := ctx.ShouldBindJSON(&s.dto); err != nil {
		return model.User{}, err
	}
	var user model.User
	s.dto.Transfer(&user)
	if err := s.repository.Create(&user); err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (s *userService) IsUnique(ctx *gin.Context) (bool, error) {
	result, err := s.repository.IsUnique(ctx.Query("email"))
	return result, err
}

func (s *userService) Activate(ctx *gin.Context) error {
	tokenString := ctx.Query("token")
	claim, err := s.jwtService.VerifyJWT(tokenString)
	if err != nil {
		return err
	}

	user, err := s.repository.Find(claim.ID)
	if err != nil {
		return err
	}

	if user.Activated {
		return config.AlreadyActivatedUserError
	}

	return s.repository.Activate(&user)
}

func (c *userService) Destroy(ctx *gin.Context) error {
	currentUser := ctx.MustGet("currentUser").(model.User)
	return c.repository.Destroy(&currentUser)
}

// test用
func TestNewUserService(jwtService JWTService, r repository.UserRepository) UserService {
	return &userService{
		jwtService: jwtService,
		repository: r,
		dto:        dto.User{},
	}
}
