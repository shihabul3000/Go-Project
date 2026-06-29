package service

import (
	"context"

	"github.com/shihabul3000/Go-Project/apperrors"
	"github.com/shihabul3000/Go-Project/dto"
	"github.com/shihabul3000/Go-Project/models"
	"github.com/shihabul3000/Go-Project/repository"
)

type ZoneService interface {
	CreateZone(ctx context.Context, req dto.CreateZoneRequest) (*dto.ParkingZoneResponse, error)
	ListZones(ctx context.Context) ([]dto.ParkingZoneAvailabilityResponse, error)
	GetZone(ctx context.Context, id uint) (*dto.ParkingZoneAvailabilityResponse, error)
	UpdateZone(ctx context.Context, id uint, req dto.UpdateZoneRequest) (*dto.ParkingZoneResponse, error)
	DeleteZone(ctx context.Context, id uint) error
}

type zoneService struct {
	zones repository.ZoneRepository
}

func NewZoneService(zones repository.ZoneRepository) ZoneService {
	return &zoneService{zones: zones}
}

func (s *zoneService) CreateZone(ctx context.Context, req dto.CreateZoneRequest) (*dto.ParkingZoneResponse, error) {
	zone := &models.ParkingZone{
		Name:          req.Name,
		Type:          req.Type,
		TotalCapacity: req.TotalCapacity,
		PricePerHour:  req.PricePerHour,
	}

	if err := s.zones.Create(ctx, zone); err != nil {
		return nil, err
	}

	response := toParkingZoneResponse(*zone)
	return &response, nil
}

func (s *zoneService) ListZones(ctx context.Context) ([]dto.ParkingZoneAvailabilityResponse, error) {
	zones, err := s.zones.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]dto.ParkingZoneAvailabilityResponse, 0, len(zones))
	for _, zone := range zones {
		activeCount, err := s.zones.CountActiveReservations(ctx, zone.ID)
		if err != nil {
			return nil, err
		}
		response = append(response, toParkingZoneAvailabilityResponse(zone, activeCount))
	}

	return response, nil
}

func (s *zoneService) GetZone(ctx context.Context, id uint) (*dto.ParkingZoneAvailabilityResponse, error) {
	zone, err := s.zones.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	activeCount, err := s.zones.CountActiveReservations(ctx, zone.ID)
	if err != nil {
		return nil, err
	}

	response := toParkingZoneAvailabilityResponse(*zone, activeCount)
	return &response, nil
}

func (s *zoneService) UpdateZone(ctx context.Context, id uint, req dto.UpdateZoneRequest) (*dto.ParkingZoneResponse, error) {
	zone, err := s.zones.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.TotalCapacity != nil {
		activeCount, err := s.zones.CountActiveReservations(ctx, id)
		if err != nil {
			return nil, err
		}
		if int64(*req.TotalCapacity) < activeCount {
			return nil, apperrors.ErrCapacityBelowActive
		}
		zone.TotalCapacity = *req.TotalCapacity
	}
	if req.Name != nil {
		zone.Name = *req.Name
	}
	if req.Type != nil {
		zone.Type = *req.Type
	}
	if req.PricePerHour != nil {
		zone.PricePerHour = *req.PricePerHour
	}

	if err := s.zones.Update(ctx, zone); err != nil {
		return nil, err
	}

	response := toParkingZoneResponse(*zone)
	return &response, nil
}

func (s *zoneService) DeleteZone(ctx context.Context, id uint) error {
	activeCount, err := s.zones.CountActiveReservations(ctx, id)
	if err != nil {
		return err
	}
	if activeCount > 0 {
		return apperrors.ErrZoneHasActiveReservations
	}

	return s.zones.Delete(ctx, id)
}

func toParkingZoneResponse(zone models.ParkingZone) dto.ParkingZoneResponse {
	return dto.ParkingZoneResponse{
		ID:            zone.ID,
		Name:          zone.Name,
		Type:          zone.Type,
		TotalCapacity: zone.TotalCapacity,
		PricePerHour:  zone.PricePerHour,
		CreatedAt:     zone.CreatedAt,
		UpdatedAt:     zone.UpdatedAt,
	}
}

func toParkingZoneAvailabilityResponse(zone models.ParkingZone, activeCount int64) dto.ParkingZoneAvailabilityResponse {
	available := zone.TotalCapacity - int(activeCount)
	if available < 0 {
		available = 0
	}

	return dto.ParkingZoneAvailabilityResponse{
		ID:             zone.ID,
		Name:           zone.Name,
		Type:           zone.Type,
		TotalCapacity:  zone.TotalCapacity,
		AvailableSpots: available,
		PricePerHour:   zone.PricePerHour,
		CreatedAt:      zone.CreatedAt,
		UpdatedAt:      zone.UpdatedAt,
	}
}
