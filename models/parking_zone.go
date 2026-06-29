package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	ZoneTypeGeneral    = "general"
	ZoneTypeEVCharging = "ev_charging"
	ZoneTypeCovered    = "covered"
)

type ParkingZone struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Name          string         `gorm:"type:varchar(140);not null" json:"name"`
	Type          string         `gorm:"type:varchar(30);not null;check:type IN ('general','ev_charging','covered')" json:"type"`
	TotalCapacity int            `gorm:"not null;check:total_capacity > 0" json:"total_capacity"`
	PricePerHour  float64        `gorm:"type:numeric(10,2);not null;check:price_per_hour > 0" json:"price_per_hour"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	Reservations []Reservation `gorm:"foreignKey:ZoneID" json:"-"`
}
