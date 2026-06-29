package apperrors

import (
	"errors"
	"net/http"
)

var (
	ErrInvalidRequest            = errors.New("invalid request body")
	ErrInvalidID                 = errors.New("invalid id parameter")
	ErrInvalidCredentials        = errors.New("invalid email or password")
	ErrUnauthorized              = errors.New("unauthorized")
	ErrForbidden                 = errors.New("forbidden")
	ErrNotFound                  = errors.New("resource not found")
	ErrDuplicateEmail            = errors.New("email already exists")
	ErrZoneFull                  = errors.New("parking zone is full")
	ErrZoneHasActiveReservations = errors.New("parking zone has active reservations")
	ErrCapacityBelowActive       = errors.New("capacity cannot be lower than active reservations")
	ErrReservationCancelled      = errors.New("reservation is already cancelled")
)

type ValidationError struct {
	Fields map[string]string
}

func (e *ValidationError) Error() string {
	return "validation failed"
}

func ToHTTP(err error) (int, string, interface{}) {
	if err == nil {
		return http.StatusInternalServerError, "Unexpected server error", nil
	}

	var validationErr *ValidationError
	if errors.As(err, &validationErr) {
		return http.StatusBadRequest, "Validation failed", validationErr.Fields
	}

	switch {
	case errors.Is(err, ErrInvalidRequest), errors.Is(err, ErrInvalidID), errors.Is(err, ErrDuplicateEmail):
		return http.StatusBadRequest, err.Error(), nil
	case errors.Is(err, ErrInvalidCredentials), errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized, err.Error(), nil
	case errors.Is(err, ErrForbidden):
		return http.StatusForbidden, err.Error(), nil
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound, err.Error(), nil
	case errors.Is(err, ErrZoneFull), errors.Is(err, ErrZoneHasActiveReservations), errors.Is(err, ErrCapacityBelowActive), errors.Is(err, ErrReservationCancelled):
		return http.StatusConflict, err.Error(), nil
	default:
		return http.StatusInternalServerError, "Unexpected server error", nil
	}
}
