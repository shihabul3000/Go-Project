package dto

import "time"

type CreateReservationRequest struct {
	ZoneID       uint   `json:"zone_id" validate:"required,gt=0"`
	LicensePlate string `json:"license_plate" validate:"required,max=15"`
}

type ReservationResponse struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	ZoneID       uint      `json:"zone_id"`
	LicensePlate string    `json:"license_plate"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type MyReservationResponse struct {
	ID           uint                    `json:"id"`
	LicensePlate string                  `json:"license_plate"`
	Status       string                  `json:"status"`
	Zone         ReservationZoneResponse `json:"zone"`
	CreatedAt    time.Time               `json:"created_at"`
}

type ReservationUserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type AdminReservationResponse struct {
	ID           uint                    `json:"id"`
	UserID       uint                    `json:"user_id"`
	ZoneID       uint                    `json:"zone_id"`
	LicensePlate string                  `json:"license_plate"`
	Status       string                  `json:"status"`
	User         ReservationUserResponse `json:"user"`
	Zone         ReservationZoneResponse `json:"zone"`
	CreatedAt    time.Time               `json:"created_at"`
	UpdatedAt    time.Time               `json:"updated_at"`
}
