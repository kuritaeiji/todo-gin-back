package controller_test

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/kuritaeiji/todo-gin-back/controller"
	"github.com/kuritaeiji/todo-gin-back/mock_service"
	"github.com/kuritaeiji/todo-gin-back/model"
	"github.com/stretchr/testify/assert"
)

var (
	assertion *assert.Assertions

	ctrl controller.UserController

	userServiceMock  *mock_service.MockUserService
	emailServiceMock *mock_service.MockEmailService

	ctx *gin.Context
	rec *httptest.ResponseRecorder
)

func prepareTest(t *testing.T) {
	assertion = assert.New(t)

	userMockCtrl := gomock.NewController(t)
	emailMockCtrl := gomock.NewController(t)
	userServiceMock = mock_service.NewMockUserService(userMockCtrl)
	emailServiceMock = mock_service.NewMockEmailService(emailMockCtrl)

	ctrl = controller.TestNewUserController(userServiceMock, emailServiceMock)

	rec = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(rec)
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	m.Run()
}

func TestCreate(t *testing.T) {
	prepareTest(t)
	user := model.User{Email: "user@example.com", PasswordDigest: "password"}
	userServiceMock.EXPECT().Create(ctx).Return(user, nil)
	emailServiceMock.EXPECT().ActivationUserEmail(user)

	ctrl.Create(ctx)
	assertion.Equal(200, rec.Code)
}

func TestInvalidCreate(t *testing.T) {
	prepareTest(t)
	var verr validator.ValidationErrors
	var err error = verr
	userServiceMock.EXPECT().Create(ctx).Return(model.User{}, err)

	ctrl.Create(ctx)
	assertion.Equal(400, rec.Code)
}

func TestIsUnique(t *testing.T) {
	prepareTest(t)
	userServiceMock.EXPECT().IsUnique(ctx).Return(true)

	ctrl.IsUnique(ctx)
	assertion.Equal(200, rec.Code)
}

func TestBadIsUnique(t *testing.T) {
	prepareTest(t)
	userServiceMock.EXPECT().IsUnique(ctx).Return(false)

	ctrl.IsUnique(ctx)
	assertion.Equal(400, rec.Code)
}
