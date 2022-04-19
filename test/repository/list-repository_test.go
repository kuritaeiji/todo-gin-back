package repository_test

import (
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kuritaeiji/todo-gin-back/config"
	"github.com/kuritaeiji/todo-gin-back/db"
	"github.com/kuritaeiji/todo-gin-back/factory"
	"github.com/kuritaeiji/todo-gin-back/model"
	"github.com/kuritaeiji/todo-gin-back/repository"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type ListRepositoryTestSuite struct {
	suite.Suite
	repository repository.ListRepository
}

func (suite *ListRepositoryTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	config.Init()
	db.Init()

	suite.repository = repository.NewListRepository()
}

func (suite *ListRepositoryTestSuite) TearDownTest() {
	db.DeleteAll()
}

func (suite *ListRepositoryTestSuite) TearDownSuite() {
	db.CloseDB()
}

func TestListRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(ListRepositoryTestSuite))
}

func (suite *ListRepositoryTestSuite) TestSuccessCreate() {
	user := factory.CreateUser(&factory.UserConfig{})
	list := factory.NewList(&factory.ListConfig{})
	err := suite.repository.Create(&user, &list)

	suite.Nil(err)
	suite.Equal(list.Title, user.Lists[0].Title)
}

func (suite *ListRepositoryTestSuite) TestSuccessUpdate() {
	user := factory.CreateUser(&factory.UserConfig{})
	list := factory.CreateList(&factory.ListConfig{Index: 1}, user)
	updatingList := factory.NewList(&factory.ListConfig{Title: "test title"})
	err := suite.repository.Update(&list, updatingList)

	suite.Nil(err)
	suite.Equal(updatingList.Title, list.Title)
	suite.Equal(1, list.Index)
}

func (suite *ListRepositoryTestSuite) TestSuccessDestroy() {
	user := factory.CreateUser(&factory.UserConfig{})
	lists := make([]model.List, 0, 4)
	for i := 0; i <= 3; i++ {
		lists = append(lists, factory.CreateList(&factory.ListConfig{Index: i}, user))
	}
	suite.repository.Destroy(&lists[2])

	_, err := suite.repository.Find(lists[2].ID)
	suite.Equal(gorm.ErrRecordNotFound, err)
	list3, _ := suite.repository.Find(lists[3].ID)
	suite.Equal(2, list3.Index)
}

func (suite *ListRepositoryTestSuite) TestSuccessMoveWhenIncreaseIndex() {
	user := factory.CreateUser(&factory.UserConfig{})
	lists := make([]model.List, 0, 5)
	for i := 0; i <= 4; i++ {
		lists = append(lists, factory.CreateList(&factory.ListConfig{Index: i, Title: strconv.Itoa(i)}, user))
	}
	err := suite.repository.Move(&lists[1], 3, &user)

	suite.Nil(err)
	suite.repository.FindLists(&user)
	suite.Equal("1", user.Lists[3].Title)
	suite.Equal("3", user.Lists[2].Title)
	suite.Equal("2", user.Lists[1].Title)

	suite.Equal("0", user.Lists[0].Title)
	suite.Equal("4", user.Lists[4].Title)
}

func (suite *ListRepositoryTestSuite) TestSuccessMoveWhenDecreaseIndex() {
	user := factory.CreateUser(&factory.UserConfig{})
	lists := make([]model.List, 0, 5)
	for i := 0; i <= 4; i++ {
		lists = append(lists, factory.CreateList(&factory.ListConfig{Index: i, Title: strconv.Itoa(i)}, user))
	}
	err := suite.repository.Move(&lists[3], 1, &user)

	suite.Nil(err)
	suite.repository.FindLists(&user)
	suite.Equal("0", user.Lists[0].Title)
	suite.Equal("3", user.Lists[1].Title)
	suite.Equal("1", user.Lists[2].Title)
	suite.Equal("2", user.Lists[3].Title)
	suite.Equal("4", user.Lists[4].Title)
}

func (suite *ListRepositoryTestSuite) TestSuccessFind() {
	user := factory.CreateUser(&factory.UserConfig{})
	list := factory.CreateList(&factory.ListConfig{}, user)
	rList, err := suite.repository.Find(list.ID)

	suite.Nil(err)
	suite.Equal(list.Title, rList.Title)
	suite.Equal(list.ID, rList.ID)
}

func (suite *ListRepositoryTestSuite) TestSuccessFindLists() {
	user := factory.CreateUser(&factory.UserConfig{})
	list1 := factory.CreateList(&factory.ListConfig{}, user)
	list2 := factory.CreateList(&factory.ListConfig{Index: 1}, user)
	err := suite.repository.FindLists(&user)

	suite.Nil(err)
	suite.Equal([]model.List{list1, list2}, user.Lists)
}
