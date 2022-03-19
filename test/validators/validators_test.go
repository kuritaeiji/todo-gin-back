package validators_test

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/kuritaeiji/todo-gin-back/db"
	"github.com/kuritaeiji/todo-gin-back/model"
	"github.com/kuritaeiji/todo-gin-back/validators"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var validate *validator.Validate
var database *gorm.DB

func TestMain(m *testing.M) {
	validators.Init()
	validate = validators.GetValidate()
	db.TestInit()
	database = db.GetDB()
	m.Run()
}

func TestValidator(t *testing.T) {
	assert := assert.New(t)
	tests := map[string]struct {
		target   string
		expected bool
	}{
		"normal":      {target: "Aa0", expected: true},
		"invalidChar": {target: "-", expected: false},
		"noNum":       {target: "Aa", expected: false},
		"noUpper":     {target: "a0", expected: false},
		"noLower":     {target: "A0", expected: false},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			err := validate.Var(test.target, "password")
			if test.expected {
				assert.Nil(err)
			} else {
				assert.Error(err)
			}
		})
	}
}

func TestUniqueEmail(t *testing.T) {
	assert := assert.New(t)

	email := "user@example.com"
	tests := map[string]struct {
		callback func()
		expected bool
	}{
		"normal": {func() {}, true},
		"invalid": {func() {
			database.Create(&model.User{Email: email, PasswordDigest: "pass"})
		}, false},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			test.callback()
			err := validate.Var(email, "unique_email")
			if test.expected {
				assert.Nil(err)
			} else {
				_, ok := err.(validator.ValidationErrors)
				assert.True(ok)
			}
			db.DeleteAll()
		})
	}
}
