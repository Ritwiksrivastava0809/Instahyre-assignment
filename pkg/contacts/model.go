package contacts

import (
	"spam-search/pkg/constants"
	"spam-search/pkg/users"
	"time"
)

type Contact struct {
	ID          uint       `json:"id" gorm:"primaryKey;not null;autoIncrement"`
	UserID      uint       `json:"user_id" gorm:"not null"`
	User        users.User `json:"user" gorm:"foreignKey:UserID"`
	Name        string     `json:"name" gorm:"not null"`
	PhoneNumber string     `json:"phone_number" gorm:"not null;index"`
	Email       *string    `json:"email" gorm:"nullable"`
	CreatedAt   time.Time  `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

func (Contact) TableName() string {
	return constants.ConactsTable
}
