package model_test

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kuritaeiji/todo-gin-back/factory"
	"github.com/kuritaeiji/todo-gin-back/model"
	"github.com/stretchr/testify/suite"
)

type UserModelTestSuite struct {
	suite.Suite
	model model.User
}

func (suite *UserModelTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
}

func (suite *UserModelTestSuite) SetupTest() {
	suite.model = factory.NewUser(&factory.UserConfig{})
}

func TestUserModel(t *testing.T) {
	suite.Run(t, new(UserModelTestSuite))
}

func (suite *UserModelTestSuite) TestSuccessAuthenticate() {
	suite.True(suite.model.Authenticate(factory.DefualtPassword))
}

func (suite *UserModelTestSuite) TestBadAuthenticate() {
	suite.False(suite.model.Authenticate("invalid password"))
}

func (suite *UserModelTestSuite) TestTrueHasList() {
	var user model.User
	var list model.List
	list.UserID = user.ID

	suite.True(user.HasList(list))
}

func (suite *UserModelTestSuite) TestFalseHasList() {
	var user model.User
	var list model.List
	list.UserID = user.ID + 1

	suite.False(user.HasList(list))
}
