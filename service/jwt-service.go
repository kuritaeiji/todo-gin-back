package service

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type UserClaim struct {
	Id int `json:"id"`
	jwt.StandardClaims
}

type JWTService interface {
	CreateJWT(id int, dayFromNow int) string
}

type jwtService struct {
	key []byte
}

func NewJWTService() JWTService {
	return &jwtService{[]byte(os.Getenv("JWT_SECRET_KEY"))}
}

func (s *jwtService) CreateJWT(id int, dayFromNow int) string {
	claim := UserClaim{
		id,
		jwt.StandardClaims{ExpiresAt: time.Now().AddDate(0, 0, dayFromNow).Unix()},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString(s.key)
	if err != nil {
		panic(err)
	}
	return tokenString
}
