package users

import (
	"spam-search/pkg/constants"
	"time"
)

type User struct {
	ID           uint      `json:"id" gorm:"primaryKey;not null;autoIncrement"`
	Name         string    `json:"name" gorm:"not null"`
	PhoneNumber  string    `json:"phone_number" gorm:"not null"`
	Email        string    `json:"email" gorm:"not null"`
	PasswordHash string    `json:"password_hash" gorm:"not null"` // Updated field name for convention
	CreatedAt    time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
}

func (User) TableName() string {
	return constants.UserTable
}