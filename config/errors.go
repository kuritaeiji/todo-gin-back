package config

import (
	"errors"

	"github.com/gin-gonic/gin"
)

var UniqueUserError = errors.New("unique user error")
var AlreadyActivatedUserError = errors.New("alreay activated user")

type ErrorResponse struct {
	Code int
	Json gin.H
}

var (
	JWTExpiredErrorResponse = ErrorResponse{
		Code: 401,
		Json: createJson("jwt expired error"),
	}

	JWTValidationErrorResponse = ErrorResponse{
		Code: 401,
		Json: createJson("jwt validation error"),
	}

	RecordNotFoundErrorResponse = ErrorResponse{
		Code: 404,
		Json: createJson("record not found"),
	}

	AlreadyActivatedUserErrorResponse = ErrorResponse{
		Code: 401,
		Json: createJson("already activated user"),
	}

	UniqueUserErrorResponse = ErrorResponse{
		Code: 400,
		Json: createJson("not unique user"),
	}

	ValidationErrorReesponse = ErrorResponse{
		Code: 400,
		Json: createJson("validation error"),
	}
)

func createJson(content string) gin.H {
	return gin.H{
		"content": content,
	}
}
