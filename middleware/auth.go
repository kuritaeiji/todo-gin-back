package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/kuritaeiji/todo-gin-back/config"
	"github.com/kuritaeiji/todo-gin-back/repository"
	"github.com/kuritaeiji/todo-gin-back/service"
)

type AuthMiddleware interface {
	Auth(*gin.Context)
	Guest(*gin.Context)
}

type authMiddleware struct {
	jwtService     service.JWTService
	userRepository repository.UserRepository
}

func NewAuthMiddleware() AuthMiddleware {
	return &authMiddleware{
		jwtService:     service.NewJWTService(),
		userRepository: repository.NewUserRepository(),
	}
}

func (m *authMiddleware) Auth(ctx *gin.Context) {
	claim, err := m.jwtService.VerifyJWT(m.tokenString(ctx))
	verr, ok := err.(*jwt.ValidationError)
	if ok && verr.Errors == jwt.ValidationErrorExpired {
		ctx.AbortWithStatusJSON(config.NotLoggedInWithJwtIsExpiredErrorResponse.Code, config.NotLoggedInWithJwtIsExpiredErrorResponse.Json)
		return
	}

	if err != nil {
		ctx.AbortWithStatusJSON(config.NotLoggedInErrorResponse.Code, config.NotLoggedInErrorResponse.Json)
		return
	}

	currentUser, err := m.userRepository.Find(claim.ID)

	if err != nil {
		ctx.AbortWithStatusJSON(config.NotLoggedInErrorResponse.Code, config.NotLoggedInErrorResponse.Json)
		return
	}

	ctx.Set(config.CurrentUserKey, currentUser)
	ctx.Next()
}

func (m *authMiddleware) Guest(ctx *gin.Context) {
	if len(m.tokenString(ctx)) > 0 {
		ctx.AbortWithStatusJSON(config.GuestErrorResponse.Code, config.GuestErrorResponse.Json)
		return
	}

	ctx.Next()
}

func (m *authMiddleware) tokenString(ctx *gin.Context) string {
	return strings.Replace(ctx.GetHeader(config.TokenHeader), config.Bearer, "", 1)
}

// test
func TestNewAuthMiddleware(jwtService service.JWTService, userRepository repository.UserRepository) AuthMiddleware {
	return &authMiddleware{
		jwtService:     jwtService,
		userRepository: userRepository,
	}
}
