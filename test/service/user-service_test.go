package service_test

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/kuritaeiji/todo-gin-back/mock_repository"
	"github.com/kuritaeiji/todo-gin-back/model"
	"github.com/kuritaeiji/todo-gin-back/service"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

var (
	userService        service.UserService
	userRepositoryMock *mock_repository.MockUserRepository
)

func prepareUserService(t *testing.T) {
	assertion = assert.New(t)
	ctrl := gomock.NewController(t)
	userRepositoryMock = mock_repository.NewMockUserRepository(ctrl)
	userService = service.TestNewUserService(userRepositoryMock)
	rec = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(rec)
}

func TestCreate(t *testing.T) {
	email := "user@example.com"
	password := "Password1010"
	body := fmt.Sprintf(`{"email":"%v","password":"%v"}`, email, password)

	t.Run("success", func(t *testing.T) {
		prepareUserService(t)
		bodyReader := strings.NewReader(body)
		req := httptest.NewRequest("POST", "/users/unique", bodyReader)
		req.Header.Add("Content-Type", binding.MIMEJSON)
		ctx.Request = req

		userRepositoryMock.EXPECT().Create(gomock.Any()).Return(nil).Do(func(actualUser *model.User) {
			assertion.Equal(actualUser.Email, email)
			assertion.Nil(bcrypt.CompareHashAndPassword([]byte(actualUser.PasswordDigest), []byte(password)))
		})
		user, err := userService.Create(ctx)

		assertion.Equal(user.Email, email)
		assertion.IsType(model.User{}, user)
		assertion.Nil(err)
	})

	t.Run("bad validation error", func(t *testing.T) {
		prepareUserService(t)
		body := strings.NewReader(`{"email":"","password":""}`)
		req := httptest.NewRequest("POST", "/users/unique", body)
		req.Header.Add("Content-Type", binding.MIMEJSON)
		ctx.Request = req

		user, err := userService.Create(ctx)

		assertion.Equal(model.User{}, user)
		assertion.IsType(validator.ValidationErrors{}, err)
	})

	t.Run("bad db error", func(t *testing.T) {
		prepareUserService(t)
		bodyReader := strings.NewReader(body)
		req := httptest.NewRequest("POST", "/users/unique", bodyReader)
		req.Header.Add("Content-Type", binding.MIMEJSON)
		ctx.Request = req

		err := errors.New("db error")
		userRepositoryMock.EXPECT().Create(gomock.Any()).Return(err).Do(func(actualUser *model.User) {
			assertion.Equal(actualUser.Email, email)
			assertion.Nil(bcrypt.CompareHashAndPassword([]byte(actualUser.PasswordDigest), []byte(password)))
		})

		user, returnErr := userService.Create(ctx)

		assertion.Equal(model.User{}, user)
		assertion.Equal(err, returnErr)
	})
}

func TestIsUnique(t *testing.T) {
	email := "user@example.com"
	tests := []struct {
		name        string
		prepareMock func(*testing.T)
		result      bool
	}{
		{
			name: "unique", prepareMock: func(t *testing.T) {
				userRepositoryMock.EXPECT().IsUnique(email).Return(true, nil)
			}, result: true,
		},
		{
			name: "not unique", prepareMock: func(t *testing.T) {
				userRepositoryMock.EXPECT().IsUnique(email).Return(false, nil)
			}, result: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prepareUserService(t)
			tt.prepareMock(t)
			req := httptest.NewRequest("GET", fmt.Sprintf("/users/unique?email=%v", email), nil)
			ctx.Request = req
			result := userService.IsUnique(ctx)

			assertion.Equal(tt.result, result)
		})
	}
}
