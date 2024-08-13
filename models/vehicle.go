package models

import (
	"time"

	"gorm.io/gorm"
)

type Vehicle struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	Make             string         `json:"make"`
	Model            string         `json:"model"`
	LicensePlate     string         `json:"license_plate"`
	AssignedDriverID *uint          `json:"assigned_driver_id,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
