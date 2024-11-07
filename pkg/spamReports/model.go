package spamreports

import (
	"spam-search/pkg/constants"
	"spam-search/pkg/users"
	"time"
)

type SpamReport struct {
	ID          uint       `json:"id" gorm:"primaryKey;not null;autoIncrement"`
	PhoneNumber string     `json:"phone_number" gorm:"not null;index"`
	ReportedBy  uint       `json:"reported_by" gorm:"not null"`
	User        users.User `json:"user" gorm:"foreignKey:ReportedBy;references:ID"`
	ReportCount uint       `json:"report_count" gorm:"default:1"`
	ReportedAt  time.Time  `json:"reported_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	Status      string     `json:"status" gorm:"type:varchar(20);default:'pending'"`
}

func (SpamReport) TableName() string {
	return constants.SpamReport
}
