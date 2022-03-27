package config

import (
	"errors"

	"github.com/gin-gonic/gin"
)

var (
	UniqueUserError             = errors.New("not unique user")
	AlreadyActivatedUserError   = errors.New("alreay activated user")
	PasswordAuthenticationError = errors.New("password is not authenticated")
	EmailClientError            = errors.New("email client error")
)

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
		Json: createJson(AlreadyActivatedUserError.Error()),
	}

	UniqueUserErrorResponse = ErrorResponse{
		Code: 400,
		Json: createJson(UniqueUserError.Error()),
	}

	ValidationErrorResponse = ErrorResponse{
		Code: 400,
		Json: createJson("validation error"),
	}

	EmailClientErrorResponse = ErrorResponse{
		Code: 500,
		Json: createJson(EmailClientError.Error()),
	}

	PasswordAuthenticationErrorResponse = ErrorResponse{
		Code: 401,
		Json: createJson(PasswordAuthenticationError.Error()),
	}

	GuestErrorResponse = ErrorResponse{
		Code: 401,
		Json: createJson("user is already logged in"),
	}

	NotLoggedInErrorResponse = ErrorResponse{
		Code: 401,
		Json: createJson("user is not logged in"),
	}
)

func createJson(content string) gin.H {
	return gin.H{
		"content": content,
	}
}
