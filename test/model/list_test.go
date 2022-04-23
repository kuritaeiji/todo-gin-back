package model_test

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kuritaeiji/todo-gin-back/factory"
	"github.com/kuritaeiji/todo-gin-back/model"
	"github.com/stretchr/testify/suite"
)

type ListModelTestSuite struct {
	suite.Suite
}

func (suite *ListModelTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
}

func TestListModelTestSuite(t *testing.T) {
	suite.Run(t, new(ListModelTestSuite))
}

func (suite *ListModelTestSuite) TestToJsonMethod() {
	list := factory.NewList(&factory.ListConfig{Title: "title"})
	json := list.ToJson()
	suite.Equal(gin.H{"title": list.Title, "id": list.ID}, json)
}

func (suite *ListModelTestSuite) TestToJsonListSlice() {
	lists := make([]model.List, 0, 2)
	for i := 0; i <= 1; i++ {
		lists = append(lists, factory.NewList(&factory.ListConfig{}))
	}
	listsJsonSlice := model.ToJsonListSlice(lists)
	suite.Equal([]gin.H{{"id": lists[0].ID, "title": lists[0].Title}, {"id": lists[1].ID, "title": lists[1].Title}}, listsJsonSlice)
}
