package test

import (
	"testing"

	"github.com/kuritaeiji/todo-gin-back/db"
	"github.com/stretchr/testify/assert"
)

func add(a, b int) int {
	return a + b
}

func TestMain(m *testing.M) {
	db.TestInit()
	defer db.CloseDB()
	m.Run()
}

func TestAdd(t *testing.T) {
	assert := assert.New(t)

	n := add(1, 2)
	assert.Equal(3, n)
}
