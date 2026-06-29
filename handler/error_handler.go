package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shihabul3000/Go-Project/apperrors"
)

func HTTPErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	if echoErr, ok := err.(*echo.HTTPError); ok {
		switch echoErr.Code {
		case http.StatusNotFound:
			err = apperrors.ErrNotFound
		case http.StatusUnauthorized:
			err = apperrors.ErrUnauthorized
		case http.StatusForbidden:
			err = apperrors.ErrForbidden
		case http.StatusBadRequest:
			err = apperrors.ErrInvalidRequest
		default:
			err = nil
		}
	}

	_ = fail(c, err)
}
