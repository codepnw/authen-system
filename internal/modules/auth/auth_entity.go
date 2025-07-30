package auth

import (
	"time"
)

type RefreshToken struct {
	ID           int64     `json:"id" gorm:"primaryKey"`
	UserID       int64     `json:"user_id" gorm:"not null"`
	RefreshToken string    `json:"token" gorm:"not null"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}
