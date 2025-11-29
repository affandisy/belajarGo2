package user

import (
	"belajarGo2/service/user"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Controller struct {
	logger  *slog.Logger
	userSvc user.Service
}

func NewController(logger *slog.Logger, s user.Service) *Controller {
	return &Controller{
		logger:  logger,
		userSvc: s,
	}
}

type userRegisterRequest struct {
	Email    string `json:"emai" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Fullname string `json:"fullname" validate:"required"`
}

func (ctrl *Controller) Register(c echo.Context) error {
	request := new(userRegisterRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": http.StatusText(http.StatusBadRequest)})
	}

	if err := validator.New().Struct(request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": http.StatusText(http.StatusBadRequest)})
	}

	_, err := ctrl.userSvc.Register(user.User{
		Email:    request.Email,
		Password: request.Password,
		Fullname: request.Fullname,
	})
	if err != nil {
		if strings.Contains(err.Error(), "registered") {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": http.StatusText(http.StatusBadRequest)})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": http.StatusText(http.StatusInternalServerError)})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{"message": http.StatusText(http.StatusCreated)})
}

type userLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (ctrl *Controller) Login(c echo.Context) error {
	request := new(userLoginRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "OK", "data": "data"})
	}

	if err := validator.New().Struct(request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": http.StatusText(http.StatusBadRequest)})
	}

	accessToken, err := ctrl.userSvc.Login(request.Email, request.Password)
	if err != nil {
		if strings.Contains(err.Error(), "wrong email") {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"message": http.StatusText(http.StatusUnauthorized)})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": http.StatusText(http.StatusInternalServerError)})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "OK", "data": accessToken})
}
