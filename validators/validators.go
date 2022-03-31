package validators

import (
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func Init() {
	var ok bool
	if validate, ok = binding.Validator.Engine().(*validator.Validate); ok {
		validate.RegisterValidation("password", password)
	}
}

func GetValidate() *validator.Validate {
	return validate
}

func password(fl validator.FieldLevel) bool {
	re1, _ := regexp.Compile(`[a-zA-Z0-9]+`)
	re2, _ := regexp.Compile(`[a-z]+`)
	re3, _ := regexp.Compile(`[A-Z]+`)
	re4, _ := regexp.Compile(`[0-9]+`)

	value := fl.Field().String()
	if re1.MatchString(value) && re2.MatchString(value) && re3.MatchString(value) && re4.MatchString(value) {
		return true
	}
	return false
}
