package service

import (
	"context"
	"strings"

	"github.com/shihabul3000/Go-Project/apperrors"
	"github.com/shihabul3000/Go-Project/dto"
	"github.com/shihabul3000/Go-Project/models"
	"github.com/shihabul3000/Go-Project/repository"
)

type ReservationService interface {
	CreateReservation(ctx context.Context, userID uint, req dto.CreateReservationRequest) (*dto.ReservationResponse, error)
	GetMyReservations(ctx context.Context, userID uint) ([]dto.MyReservationResponse, error)
	CancelReservation(ctx context.Context, reservationID uint, requesterID uint, requesterRole string) error
	GetAllReservations(ctx context.Context) ([]dto.AdminReservationResponse, error)
}

type reservationService struct {
	reservations repository.ReservationRepository
}

func NewReservationService(reservations repository.ReservationRepository) ReservationService {
	return &reservationService{reservations: reservations}
}

func (s *reservationService) CreateReservation(ctx context.Context, userID uint, req dto.CreateReservationRequest) (*dto.ReservationResponse, error) {
	licensePlate := strings.ToUpper(strings.TrimSpace(req.LicensePlate))
	reservation, err := s.reservations.CreateWithCapacityLock(ctx, userID, req.ZoneID, licensePlate)
	if err != nil {
		return nil, err
	}

	response := toReservationResponse(*reservation)
	return &response, nil
}

func (s *reservationService) GetMyReservations(ctx context.Context, userID uint) ([]dto.MyReservationResponse, error) {
	reservations, err := s.reservations.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	response := make([]dto.MyReservationResponse, 0, len(reservations))
	for _, reservation := range reservations {
		response = append(response, dto.MyReservationResponse{
			ID:           reservation.ID,
			LicensePlate: reservation.LicensePlate,
			Status:       reservation.Status,
			Zone:         toReservationZoneResponse(reservation.Zone),
			CreatedAt:    reservation.CreatedAt,
		})
	}

	return response, nil
}

func (s *reservationService) CancelReservation(ctx context.Context, reservationID uint, requesterID uint, requesterRole string) error {
	reservation, err := s.reservations.FindByID(ctx, reservationID)
	if err != nil {
		return err
	}

	if requesterRole != models.RoleAdmin && reservation.UserID != requesterID {
		return apperrors.ErrForbidden
	}

	if reservation.Status == models.ReservationStatusCancelled {
		return apperrors.ErrReservationCancelled
	}

	reservation.Status = models.ReservationStatusCancelled
	return s.reservations.Update(ctx, reservation)
}

func (s *reservationService) GetAllReservations(ctx context.Context) ([]dto.AdminReservationResponse, error) {
	reservations, err := s.reservations.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]dto.AdminReservationResponse, 0, len(reservations))
	for _, reservation := range reservations {
		response = append(response, dto.AdminReservationResponse{
			ID:           reservation.ID,
			UserID:       reservation.UserID,
			ZoneID:       reservation.ZoneID,
			LicensePlate: reservation.LicensePlate,
			Status:       reservation.Status,
			User: dto.ReservationUserResponse{
				ID:    reservation.User.ID,
				Name:  reservation.User.Name,
				Email: reservation.User.Email,
				Role:  reservation.User.Role,
			},
			Zone:      toReservationZoneResponse(reservation.Zone),
			CreatedAt: reservation.CreatedAt,
			UpdatedAt: reservation.UpdatedAt,
		})
	}

	return response, nil
}

func toReservationResponse(reservation models.Reservation) dto.ReservationResponse {
	return dto.ReservationResponse{
		ID:           reservation.ID,
		UserID:       reservation.UserID,
		ZoneID:       reservation.ZoneID,
		LicensePlate: reservation.LicensePlate,
		Status:       reservation.Status,
		CreatedAt:    reservation.CreatedAt,
		UpdatedAt:    reservation.UpdatedAt,
	}
}

func toReservationZoneResponse(zone models.ParkingZone) dto.ReservationZoneResponse {
	return dto.ReservationZoneResponse{
		ID:   zone.ID,
		Name: zone.Name,
		Type: zone.Type,
	}
}
