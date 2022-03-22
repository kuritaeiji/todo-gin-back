package service_test

import (
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/kuritaeiji/todo-gin-back/service"
	"github.com/stretchr/testify/assert"
)

var jwtService service.JWTService

func prepareJWTService(t *testing.T) {
	assertion = assert.New(t)
	jwtService = service.NewJWTService()
	rec = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(rec)
}

func TestCreateJWT(t *testing.T) {
	prepareJWTService(t)
	id := 1
	dayFromNow := 1
	tokenString := jwtService.CreateJWT(id, dayFromNow)

	token, _ := jwt.ParseWithClaims(tokenString, &service.UserClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	claim, _ := token.Claims.(*service.UserClaim)
	assertion.Equal(1, claim.Id)
	assertion.InEpsilon(time.Now().AddDate(0, 0, dayFromNow).Unix(), claim.ExpiresAt, 10)
}
