package users

import (
	"fmt"
	"spam-search/pkg/constants"
	errorlogs "spam-search/pkg/constants/errorlogs"

	"gorm.io/gorm"
)

func (u *User) CreateUser(db *gorm.DB) error {
	tx := db.Table(constants.UserTable).Create(u)
	if tx.Error != nil {
		return fmt.Errorf(errorlogs.InsertTxnError, constants.UserTable, tx.Error)
	}
	return nil
}

// New function to check if user already exists
func (u *User) UserExists(db *gorm.DB) (bool, error) {
	var existingUser User
	// Check if a user exists with the same email or phone number
	err := db.Where("email = ? OR phone_number = ?", u.Email, u.PhoneNumber).First(&existingUser).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		// Some database error occurred
		return false, fmt.Errorf("error checking user existence: %v", err)
	}
	// If err is nil, a user was found, hence user exists
	return err == nil, nil
}
