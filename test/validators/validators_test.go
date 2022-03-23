package validators_test

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/kuritaeiji/todo-gin-back/validators"
	"github.com/stretchr/testify/assert"
)

var validate *validator.Validate

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	validators.Init()
	validate = validators.GetValidate()
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
