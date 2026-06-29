package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shihabul3000/Go-Project/apperrors"
	"github.com/shihabul3000/Go-Project/dto"
	"github.com/shihabul3000/Go-Project/service"
)

type ZoneHandler struct {
	zones service.ZoneService
}

func NewZoneHandler(zones service.ZoneService) *ZoneHandler {
	return &ZoneHandler{zones: zones}
}

func (h *ZoneHandler) Create(c echo.Context) error {
	var req dto.CreateZoneRequest
	if err := c.Bind(&req); err != nil {
		return fail(c, apperrors.ErrInvalidRequest)
	}
	if err := c.Validate(&req); err != nil {
		return fail(c, err)
	}

	zone, err := h.zones.CreateZone(c.Request().Context(), req)
	if err != nil {
		return fail(c, err)
	}

	return success(c, http.StatusCreated, "Parking zone created successfully", zone)
}

func (h *ZoneHandler) List(c echo.Context) error {
	zones, err := h.zones.ListZones(c.Request().Context())
	if err != nil {
		return fail(c, err)
	}

	return success(c, http.StatusOK, "Parking zones retrieved successfully", zones)
}

func (h *ZoneHandler) Get(c echo.Context) error {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return fail(c, err)
	}

	zone, err := h.zones.GetZone(c.Request().Context(), id)
	if err != nil {
		return fail(c, err)
	}

	return success(c, http.StatusOK, "Parking zone retrieved successfully", zone)
}

func (h *ZoneHandler) Update(c echo.Context) error {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return fail(c, err)
	}

	var req dto.UpdateZoneRequest
	if err := c.Bind(&req); err != nil {
		return fail(c, apperrors.ErrInvalidRequest)
	}
	if err := c.Validate(&req); err != nil {
		return fail(c, err)
	}

	zone, err := h.zones.UpdateZone(c.Request().Context(), id, req)
	if err != nil {
		return fail(c, err)
	}

	return success(c, http.StatusOK, "Parking zone updated successfully", zone)
}

func (h *ZoneHandler) Delete(c echo.Context) error {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return fail(c, err)
	}

	if err := h.zones.DeleteZone(c.Request().Context(), id); err != nil {
		return fail(c, err)
	}

	return successNoData(c, http.StatusOK, "Parking zone deleted successfully")
}
