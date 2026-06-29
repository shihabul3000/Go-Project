package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shihabul3000/Go-Project/apperrors"
	"github.com/shihabul3000/Go-Project/dto"
	"github.com/shihabul3000/Go-Project/service"
)

type AuthHandler struct {
	auth service.AuthService
}

func NewAuthHandler(auth service.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return fail(c, apperrors.ErrInvalidRequest)
	}
	if err := c.Validate(&req); err != nil {
		return fail(c, err)
	}

	user, err := h.auth.Register(c.Request().Context(), req)
	if err != nil {
		return fail(c, err)
	}

	return success(c, http.StatusCreated, "User registered successfully", user)
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return fail(c, apperrors.ErrInvalidRequest)
	}
	if err := c.Validate(&req); err != nil {
		return fail(c, err)
	}

	response, err := h.auth.Login(c.Request().Context(), req)
	if err != nil {
		return fail(c, err)
	}

	return success(c, http.StatusOK, "Login successful", response)
}
