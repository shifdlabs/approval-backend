package model

import (
	"time"

	"gorm.io/gorm"
)

type AppSettings struct {
	gorm.Model
	Key       string     `gorm:"type:varchar"`
	Value     string     `gorm:"type:varchar"`
	CreatedAt *time.Time `gorm:"not null;default:now()"`
	UpdatedAt *time.Time `gorm:"not null;default:now()"`
}
