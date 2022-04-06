package service

// mockgen -source=service/jwt-service.go -destination=mock_service/jwt-service.go

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/kuritaeiji/todo-gin-back/model"
)

const (
	DayFromNowAccessToken       = 365
	DayFromNowActivateUserToken = 1
)

type UserClaim struct {
	ID int `json:"id"`
	jwt.StandardClaims
}

type JWTService interface {
	CreateJWT(user model.User, dayFromNow int) string
	VerifyJWT(tokdnString string) (*UserClaim, error)
}

type jwtService struct {
	key []byte
}

func NewJWTService() JWTService {
	return &jwtService{[]byte(os.Getenv("JWT_SECRET_KEY"))}
}

func (s *jwtService) CreateJWT(user model.User, dayFromNow int) string {
	claim := UserClaim{
		user.ID,
		jwt.StandardClaims{ExpiresAt: time.Now().AddDate(0, 0, dayFromNow).Unix()},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString(s.key)
	if err != nil {
		panic(err)
	}
	return tokenString
}

func (s *jwtService) VerifyJWT(tokenString string) (*UserClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaim{}, func(t *jwt.Token) (interface{}, error) {
		return s.key, nil
	})

	if err != nil {
		return &UserClaim{}, err
	}

	if claim, ok := token.Claims.(*UserClaim); ok && token.Valid {
		return claim, nil
	}

	return &UserClaim{}, err
}
