package db

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID        string         `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primarykey"`
	CreatedAt time.Time      `json:"created_at,omitempty"`
	UpdatedAt time.Time      `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
