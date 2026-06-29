package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shihabul3000/Go-Project/apperrors"
	"github.com/shihabul3000/Go-Project/dto"
	"github.com/shihabul3000/Go-Project/middleware"
	"github.com/shihabul3000/Go-Project/service"
)

type ReservationHandler struct {
	reservations service.ReservationService
}

func NewReservationHandler(reservations service.ReservationService) *ReservationHandler {
	return &ReservationHandler{reservations: reservations}
}

func (h *ReservationHandler) Create(c echo.Context) error {
	userID, ok := middleware.UserIDFromContext(c)
	if !ok {
		return fail(c, apperrors.ErrUnauthorized)
	}

	var req dto.CreateReservationRequest
	if err := c.Bind(&req); err != nil {
		return fail(c, apperrors.ErrInvalidRequest)
	}
	if err := c.Validate(&req); err != nil {
		return fail(c, err)
	}

	reservation, err := h.reservations.CreateReservation(c.Request().Context(), userID, req)
	if err != nil {
		return fail(c, err)
	}

	return success(c, http.StatusCreated, "Reservation confirmed successfully", reservation)
}

func (h *ReservationHandler) MyReservations(c echo.Context) error {
	userID, ok := middleware.UserIDFromContext(c)
	if !ok {
		return fail(c, apperrors.ErrUnauthorized)
	}

	reservations, err := h.reservations.GetMyReservations(c.Request().Context(), userID)
	if err != nil {
		return fail(c, err)
	}

	return success(c, http.StatusOK, "My reservations retrieved successfully", reservations)
}

func (h *ReservationHandler) Cancel(c echo.Context) error {
	userID, ok := middleware.UserIDFromContext(c)
	if !ok {
		return fail(c, apperrors.ErrUnauthorized)
	}
	role, ok := middleware.RoleFromContext(c)
	if !ok {
		return fail(c, apperrors.ErrUnauthorized)
	}

	id, err := parseIDParam(c, "id")
	if err != nil {
		return fail(c, err)
	}

	if err := h.reservations.CancelReservation(c.Request().Context(), id, userID, role); err != nil {
		return fail(c, err)
	}

	return successNoData(c, http.StatusOK, "Reservation cancelled successfully")
}

func (h *ReservationHandler) ListAll(c echo.Context) error {
	reservations, err := h.reservations.GetAllReservations(c.Request().Context())
	if err != nil {
		return fail(c, err)
	}

	return success(c, http.StatusOK, "Reservations retrieved successfully", reservations)
}
