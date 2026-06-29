package models

import "time"

const (
	RoleDriver = "driver"
	RoleAdmin  = "admin"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(120);not null" json:"name"`
	Email     string    `gorm:"type:varchar(160);uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"type:varchar(255);not null" json:"-"`
	Role      string    `gorm:"type:varchar(20);not null;default:'driver';check:role IN ('driver','admin')" json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Reservations []Reservation `gorm:"foreignKey:UserID" json:"-"`
}
