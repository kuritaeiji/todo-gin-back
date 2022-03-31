package service_test

import (
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/kuritaeiji/todo-gin-back/model"
	"github.com/kuritaeiji/todo-gin-back/service"
	"github.com/stretchr/testify/suite"
)

type JWTServiceTestSuite struct {
	suite.Suite
	service service.JWTService
}

func (suite *JWTServiceTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	suite.service = service.NewJWTService()
}

func TestJWTService(t *testing.T) {
	suite.Run(t, &JWTServiceTestSuite{})
}

func (suite *JWTServiceTestSuite) TestSuccessCreateJWT() {
	id := 1
	dayFromNow := 1
	user := model.User{ID: id}
	tokenString := suite.service.CreateJWT(user, dayFromNow)
	claim, _ := suite.service.VerifyJWT(tokenString)

	suite.Equal(id, claim.ID)
	suite.InEpsilon(time.Now().AddDate(0, 0, dayFromNow).Unix(), claim.ExpiresAt, 30)
}

func (suite *JWTServiceTestSuite) TestSuccessVerifyJWT() {
	id := 1
	dayFromNow := 1
	user := model.User{ID: id}
	tokenString := suite.service.CreateJWT(user, dayFromNow)
	claim, err := suite.service.VerifyJWT(tokenString)

	suite.Equal(id, claim.ID)
	suite.Nil(err)
}

func (suite *JWTServiceTestSuite) TestBadVerifyJWTWithExpired() {
	dayFromNow := -1
	tokenString := suite.service.CreateJWT(model.User{}, dayFromNow)
	_, err := suite.service.VerifyJWT(tokenString)

	suite.IsType(&jwt.ValidationError{}, err)
}
