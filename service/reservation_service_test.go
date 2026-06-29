package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/shihabul3000/Go-Project/apperrors"
	"github.com/shihabul3000/Go-Project/dto"
	"github.com/shihabul3000/Go-Project/models"
)

type fakeReservationRepository struct {
	reservations map[uint]*models.Reservation
	nextID       uint
	createErr    error
}

func newFakeReservationRepository() *fakeReservationRepository {
	return &fakeReservationRepository{
		reservations: map[uint]*models.Reservation{},
		nextID:       1,
	}
}

func (r *fakeReservationRepository) CreateWithCapacityLock(ctx context.Context, userID uint, zoneID uint, licensePlate string) (*models.Reservation, error) {
	if r.createErr != nil {
		return nil, r.createErr
	}
	reservation := &models.Reservation{
		ID:           r.nextID,
		UserID:       userID,
		ZoneID:       zoneID,
		LicensePlate: licensePlate,
		Status:       models.ReservationStatusActive,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
	r.nextID++
	r.reservations[reservation.ID] = reservation
	return reservation, nil
}

func (r *fakeReservationRepository) FindByUserID(ctx context.Context, userID uint) ([]models.Reservation, error) {
	var reservations []models.Reservation
	for _, reservation := range r.reservations {
		if reservation.UserID == userID {
			reservations = append(reservations, *reservation)
		}
	}
	return reservations, nil
}

func (r *fakeReservationRepository) FindByID(ctx context.Context, id uint) (*models.Reservation, error) {
	reservation, ok := r.reservations[id]
	if !ok {
		return nil, apperrors.ErrNotFound
	}
	copied := *reservation
	return &copied, nil
}

func (r *fakeReservationRepository) FindAll(ctx context.Context) ([]models.Reservation, error) {
	var reservations []models.Reservation
	for _, reservation := range r.reservations {
		reservations = append(reservations, *reservation)
	}
	return reservations, nil
}

func (r *fakeReservationRepository) Update(ctx context.Context, reservation *models.Reservation) error {
	copied := *reservation
	r.reservations[reservation.ID] = &copied
	return nil
}

func TestReservationServiceCreateNormalizesLicensePlate(t *testing.T) {
	repo := newFakeReservationRepository()
	service := NewReservationService(repo)

	response, err := service.CreateReservation(context.Background(), 7, dto.CreateReservationRequest{
		ZoneID:       2,
		LicensePlate: " abc-1234 ",
	})
	if err != nil {
		t.Fatalf("CreateReservation() error = %v", err)
	}
	if response.LicensePlate != "ABC-1234" {
		t.Fatalf("expected normalized license plate, got %q", response.LicensePlate)
	}
}

func TestReservationServiceCreateRejectsBlankLicensePlate(t *testing.T) {
	repo := newFakeReservationRepository()
	service := NewReservationService(repo)

	_, err := service.CreateReservation(context.Background(), 7, dto.CreateReservationRequest{
		ZoneID:       2,
		LicensePlate: "   ",
	})

	var validationErr *apperrors.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected validation error, got %v", err)
	}
	if validationErr.Fields["license_plate"] != "is required" {
		t.Fatalf("unexpected license plate validation message: %q", validationErr.Fields["license_plate"])
	}
}

func TestReservationServicePreservesZoneFullConflict(t *testing.T) {
	repo := newFakeReservationRepository()
	repo.createErr = apperrors.ErrZoneFull
	service := NewReservationService(repo)

	_, err := service.CreateReservation(context.Background(), 7, dto.CreateReservationRequest{
		ZoneID:       2,
		LicensePlate: "ABC-1234",
	})
	if !errors.Is(err, apperrors.ErrZoneFull) {
		t.Fatalf("expected ErrZoneFull, got %v", err)
	}
}

func TestReservationServiceDriverCannotCancelAnotherDriversReservation(t *testing.T) {
	repo := newFakeReservationRepository()
	repo.reservations[10] = &models.Reservation{
		ID:           10,
		UserID:       99,
		ZoneID:       2,
		LicensePlate: "ABC-1234",
		Status:       models.ReservationStatusActive,
	}
	service := NewReservationService(repo)

	err := service.CancelReservation(context.Background(), 10, 7, models.RoleDriver)
	if !errors.Is(err, apperrors.ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestReservationServiceAdminCanCancelAnyReservation(t *testing.T) {
	repo := newFakeReservationRepository()
	repo.reservations[10] = &models.Reservation{
		ID:           10,
		UserID:       99,
		ZoneID:       2,
		LicensePlate: "ABC-1234",
		Status:       models.ReservationStatusActive,
	}
	service := NewReservationService(repo)

	if err := service.CancelReservation(context.Background(), 10, 7, models.RoleAdmin); err != nil {
		t.Fatalf("CancelReservation() error = %v", err)
	}
	if repo.reservations[10].Status != models.ReservationStatusCancelled {
		t.Fatalf("expected cancelled status, got %q", repo.reservations[10].Status)
	}
}
