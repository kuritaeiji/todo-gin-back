package repository_test

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kuritaeiji/todo-gin-back/config"
	"github.com/kuritaeiji/todo-gin-back/db"
	"github.com/kuritaeiji/todo-gin-back/model"
	"github.com/kuritaeiji/todo-gin-back/repository"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	userRepository repository.UserRepository
	db             *gorm.DB
}

func (suite *UserRepositoryTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	config.Init()
	db.Init()
	suite.userRepository = repository.NewUserRepository()
	suite.db = db.GetDB()
}

func (suite *UserRepositoryTestSuite) TearDownSuite() {
	db.CloseDB()
}

func (suite *UserRepositoryTestSuite) TearDownTest() {
	db.DeleteAll()
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}

func (suite *UserRepositoryTestSuite) TestSuccessCreate() {
	email := "user@example.com"
	password := "Password1010"

	var count int64
	suite.db.Model(&model.User{}).Count(&count)
	suite.Equal(int64(0), count)
	err := suite.userRepository.Create(&model.User{Email: email, PasswordDigest: password})
	suite.Nil(err)
	suite.db.Model(&model.User{}).Count(&count)
	suite.Equal(int64(1), count)
}

func (suite *UserRepositoryTestSuite) TestBadCreateWithNotUniqueEmail() {
	email := "user@example.com"
	user := model.User{Email: email, PasswordDigest: "password"}
	suite.db.Create(&user)

	err := suite.userRepository.Create(&user)
	suite.Equal(config.UniqueUserError, err)
}

func (suite *UserRepositoryTestSuite) TestTrueIsUnique() {
	result, _ := suite.userRepository.IsUnique("email")
	suite.True(result)
}

func (suite *UserRepositoryTestSuite) TestFalseIsUnique() {
	email := "user@example.com"
	user := model.User{Email: email, PasswordDigest: "pass"}
	suite.db.Create(&user)
	result, _ := suite.userRepository.IsUnique(email)
	suite.False(result)
}

func (suite *UserRepositoryTestSuite) TestSuccessActivate() {
	user := model.User{ID: 1}
	suite.db.Create(&user)
	err := suite.userRepository.Activate(&user)

	var ruser model.User
	suite.db.First(&ruser)
	suite.True(ruser.Activated)
	suite.Nil(err)
}

func (suite *UserRepositoryTestSuite) TestSuccessFindByEmail() {
	email := "user@example.com"
	user := model.User{Email: email}
	suite.db.Create(&user)
	ruser, err := suite.userRepository.FindByEmail(email)

	suite.Equal(user.ID, ruser.ID)
	suite.Nil(err)
}

func (suite *UserRepositoryTestSuite) TestBadFindByEmailWithRecordNotFound() {
	_, err := suite.userRepository.FindByEmail("mail")

	suite.Equal(gorm.ErrRecordNotFound, err)
}
