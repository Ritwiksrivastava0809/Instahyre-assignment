package spamreports

import (
	"spam-search/pkg/constants"
	"time"
)

type SpamReportRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	Name        string `json:"name"`
}

type GlobalSpam struct {
	ID              uint      `json:"id" gorm:"primaryKey;not null;autoIncrement"`
	PhoneNumber     string    `json:"phone_number" gorm:"not null;index"`
	Name            string    `json:"name" gorm:"type:varchar(255);index;default:null"`
	SpamReportCount uint      `json:"spam_report_count" gorm:"default:0"`
	SpamLikelihood  float64   `json:"spam_likelihood"`
	Status          string    `json:"status" gorm:"type:varchar(20);default:'pending'"`
	ReportedAt      time.Time `json:"reported_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (GlobalSpam) TableName() string {
	return constants.GlobalSpam
}
