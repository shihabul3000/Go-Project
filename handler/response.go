package handler

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/shihabul3000/Go-Project/apperrors"
	"github.com/shihabul3000/Go-Project/dto"
)

func success(c echo.Context, status int, message string, data interface{}) error {
	return c.JSON(status, dto.APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func successNoData(c echo.Context, status int, message string) error {
	return c.JSON(status, dto.APIResponse{
		Success: true,
		Message: message,
	})
}

func fail(c echo.Context, err error) error {
	status, message, details := apperrors.ToHTTP(err)
	if details == nil {
		details = message
	}
	return c.JSON(status, dto.APIResponse{
		Success: false,
		Message: message,
		Errors:  details,
	})
}

func parseIDParam(c echo.Context, name string) (uint, error) {
	value := c.Param(name)
	id, err := strconv.ParseUint(value, 10, 64)
	if err != nil || id == 0 {
		return 0, apperrors.ErrInvalidID
	}
	return uint(id), nil
}
