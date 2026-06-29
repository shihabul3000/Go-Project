package repository

import (
	"context"
	"errors"

	"github.com/shihabul3000/Go-Project/apperrors"
	"github.com/shihabul3000/Go-Project/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ReservationRepository interface {
	CreateWithCapacityLock(ctx context.Context, userID uint, zoneID uint, licensePlate string) (*models.Reservation, error)
	FindByUserID(ctx context.Context, userID uint) ([]models.Reservation, error)
	FindByID(ctx context.Context, id uint) (*models.Reservation, error)
	FindAll(ctx context.Context) ([]models.Reservation, error)
	Update(ctx context.Context, reservation *models.Reservation) error
}

type reservationRepository struct {
	db *gorm.DB
}

func NewReservationRepository(db *gorm.DB) ReservationRepository {
	return &reservationRepository{db: db}
}

func (r *reservationRepository) CreateWithCapacityLock(ctx context.Context, userID uint, zoneID uint, licensePlate string) (*models.Reservation, error) {
	var reservation *models.Reservation

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var zone models.ParkingZone
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&zone, zoneID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperrors.ErrNotFound
			}
			return err
		}

		var activeCount int64
		if err := tx.Model(&models.Reservation{}).
			Where("zone_id = ? AND status = ?", zoneID, models.ReservationStatusActive).
			Count(&activeCount).Error; err != nil {
			return err
		}

		if activeCount >= int64(zone.TotalCapacity) {
			return apperrors.ErrZoneFull
		}

		reservation = &models.Reservation{
			UserID:       userID,
			ZoneID:       zoneID,
			LicensePlate: licensePlate,
			Status:       models.ReservationStatusActive,
		}

		return tx.Create(reservation).Error
	})
	if err != nil {
		return nil, err
	}

	return reservation, nil
}

func (r *reservationRepository) FindByUserID(ctx context.Context, userID uint) ([]models.Reservation, error) {
	var reservations []models.Reservation
	err := r.db.WithContext(ctx).
		Preload("Zone", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped()
		}).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&reservations).Error
	return reservations, err
}

func (r *reservationRepository) FindByID(ctx context.Context, id uint) (*models.Reservation, error) {
	var reservation models.Reservation
	if err := r.db.WithContext(ctx).
		Preload("Zone", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped()
		}).
		Preload("User").
		First(&reservation, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return &reservation, nil
}

func (r *reservationRepository) FindAll(ctx context.Context) ([]models.Reservation, error) {
	var reservations []models.Reservation
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Zone", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped()
		}).
		Order("created_at DESC").
		Find(&reservations).Error
	return reservations, err
}

func (r *reservationRepository) Update(ctx context.Context, reservation *models.Reservation) error {
	return r.db.WithContext(ctx).Save(reservation).Error
}
