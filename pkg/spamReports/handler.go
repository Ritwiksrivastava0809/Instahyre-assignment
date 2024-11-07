package spamreports

import (
	"spam-search/pkg/constants"
	"spam-search/pkg/contacts"
	"spam-search/pkg/users"
	"time"

	"gorm.io/gorm"
)

func CalculateSpamLikelihood(reportCount uint) float64 {
	return float64(reportCount) / float64(reportCount+100)
}

func ReportSpam(db *gorm.DB, phoneNumber string, name string) (*GlobalSpam, error) {
	var spamReport GlobalSpam

	tx := db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	err := tx.Where("phone_number = ?", phoneNumber).First(&spamReport).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			spamReport = GlobalSpam{
				PhoneNumber: phoneNumber,
				Name:        name,
				Status:      constants.Pending,
				ReportedAt:  time.Now(),
			}
			if err := tx.Create(&spamReport).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		} else {
			tx.Rollback()
			return nil, err
		}
	}

	spamReport.SpamReportCount++
	spamReport.SpamLikelihood = CalculateSpamLikelihood(spamReport.SpamReportCount)

	if spamReport.SpamReportCount > 5 {
		spamReport.Status = constants.Spam
	}

	if err := tx.Save(&spamReport).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return &spamReport, nil
}

func GetSpamReportsByName(db *gorm.DB, name string) ([]GlobalSpam, error) {
	var results []GlobalSpam
	err := db.Where("name LIKE ?", name+"%").Order("name ASC").Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

func GetSpamReportsByPhoneNumber(db *gorm.DB, phoneNumber string) ([]GlobalSpam, error) {
	var results []GlobalSpam
	err := db.Where("phone_number LIKE ?", phoneNumber+"%").Order("phone_number ASC").Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

func GetUserByPhoneNumber(db *gorm.DB, phoneNumber string) (*users.User, error) {
	var user users.User
	err := db.Where("phone_number = ?", phoneNumber).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserContactEmail(db *gorm.DB, userID uint, phoneNumber string) (string, error) {
	var contact contacts.Contact
	err := db.Where("user_id = ? AND phone_number = ?", userID, phoneNumber).First(&contact).Error
	if err != nil {
		return "", err
	}
	if contact.Email != nil {
		return *contact.Email, nil
	}
	return "", nil
}
