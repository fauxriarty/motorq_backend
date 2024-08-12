package models

import (
	"time"

	"gorm.io/gorm"
)

type Assignment struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	DriverID  uint           `json:"driver_id"`
	VehicleID uint           `json:"vehicle_id"`
	StartTime time.Time      `json:"start_time"`
	EndTime   time.Time      `json:"end_time"`
	Status    string         `json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
