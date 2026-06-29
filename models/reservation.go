package models

import "time"

const (
	ReservationStatusActive    = "active"
	ReservationStatusCompleted = "completed"
	ReservationStatusCancelled = "cancelled"
)

type Reservation struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"not null;index" json:"user_id"`
	ZoneID       uint      `gorm:"not null;index" json:"zone_id"`
	LicensePlate string    `gorm:"type:varchar(15);not null;index" json:"license_plate"`
	Status       string    `gorm:"type:varchar(20);not null;default:'active';index;check:status IN ('active','completed','cancelled')" json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	User User        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"user"`
	Zone ParkingZone `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"zone"`
}
