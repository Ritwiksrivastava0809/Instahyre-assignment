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

func (u *User) UserExists(db *gorm.DB) (bool, error) {
	var existingUser User

	err := db.Where("email = ? OR phone_number = ?", u.Email, u.PhoneNumber).First(&existingUser).Error
	if err != nil && err != gorm.ErrRecordNotFound {

		return false, fmt.Errorf("error checking user existence: %v", err)
	}

	return err == nil, nil
}

func GetUserByPhoneNumber(db *gorm.DB, phoneNumber string) (User, error) {
	var user User
	err := db.Where("phone_number = ?", phoneNumber).First(&user).Error
	if err != nil {
		return user, fmt.Errorf(errorlogs.GetUserByPhoneNumberError, phoneNumber, err)
	}
	return user, nil
}
