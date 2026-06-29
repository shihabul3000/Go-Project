package repository

import (
	"context"
	"errors"

	"github.com/shihabul3000/Go-Project/apperrors"
	"github.com/shihabul3000/Go-Project/models"
	"gorm.io/gorm"
)

type ZoneRepository interface {
	Create(ctx context.Context, zone *models.ParkingZone) error
	FindAll(ctx context.Context) ([]models.ParkingZone, error)
	FindByID(ctx context.Context, id uint) (*models.ParkingZone, error)
	CountActiveReservations(ctx context.Context, zoneID uint) (int64, error)
	Update(ctx context.Context, zone *models.ParkingZone) error
	Delete(ctx context.Context, id uint) error
}

type zoneRepository struct {
	db *gorm.DB
}

func NewZoneRepository(db *gorm.DB) ZoneRepository {
	return &zoneRepository{db: db}
}

func (r *zoneRepository) Create(ctx context.Context, zone *models.ParkingZone) error {
	return r.db.WithContext(ctx).Create(zone).Error
}

func (r *zoneRepository) FindAll(ctx context.Context) ([]models.ParkingZone, error) {
	var zones []models.ParkingZone
	if err := r.db.WithContext(ctx).Order("id ASC").Find(&zones).Error; err != nil {
		return nil, err
	}
	return zones, nil
}

func (r *zoneRepository) FindByID(ctx context.Context, id uint) (*models.ParkingZone, error) {
	var zone models.ParkingZone
	if err := r.db.WithContext(ctx).First(&zone, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return &zone, nil
}

func (r *zoneRepository) CountActiveReservations(ctx context.Context, zoneID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Reservation{}).
		Where("zone_id = ? AND status = ?", zoneID, models.ReservationStatusActive).
		Count(&count).Error
	return count, err
}

func (r *zoneRepository) Update(ctx context.Context, zone *models.ParkingZone) error {
	return r.db.WithContext(ctx).Save(zone).Error
}

func (r *zoneRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.ParkingZone{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}
