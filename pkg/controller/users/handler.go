package userController

import (
	"fmt"
	"net/http"
	"net/mail"
	"spam-search/pkg/constants"
	errorlogs "spam-search/pkg/constants/errorlogs"
	"spam-search/pkg/token"
	"spam-search/pkg/users"
	"spam-search/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func (con *UserController) CreateUserHandler(c *gin.Context) {
	var user users.User

	if err := c.ShouldBindJSON(&user); err != nil {
		log.Error().Msgf(errorlogs.BindingJsonError, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": errorlogs.BindingJsonError})
		return
	}

	if user.Email != "" {
		if _, err := mail.ParseAddress(user.Email); err != nil {
			log.Error().Msg(errorlogs.InvalidEmailFormatError)
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid email format"})
			return
		}
	}

	db := c.MustGet(constants.ConstantDb).(*gorm.DB)
	if db.Error != nil {
		message := fmt.Sprintf(errorlogs.GetDBError, db.Error)
		log.Error().Msg("Error while getting the db connection" + message)
		c.AbortWithStatusJSON(http.StatusInternalServerError, message)
		return
	}

	tx := db.Begin()
	if tx.Error != nil {
		message := fmt.Sprintf(errorlogs.BeginSQLTransactionError, tx.Error)
		log.Error().Msg("Error while starting a transaction to insert user." + message)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": message})
		return
	}

	defer tx.Rollback()

	exists, err := user.UserExists(tx)
	if err != nil {
		tx.Rollback()
		log.Error().Msgf("Error checking user existence: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error checking user existence"})
		return
	}
	if exists {
		log.Error().Msg("User already exists with email or phone number")
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{"message": "User already exists with this email or phone number"})
		return
	}

	hashedPassword, err := utils.HashPasswordArgon2(user.PasswordHash)
	if err != nil {
		tx.Rollback()
		log.Error().Msgf(errorlogs.HashedPasswordError, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": errorlogs.HashedPasswordError})
		return
	}

	user.PasswordHash = hashedPassword

	if err := user.CreateUser(tx); err != nil {
		tx.Rollback()
		log.Error().Msg(errorlogs.InsertUserError + err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Commit().Error; err != nil {
		message := fmt.Sprintf(errorlogs.CommitSQLTransactionError, constants.UserTable, err)
		log.Error().Msg("Error while commiting transaction to create user. " + message)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

func (con *UserController) LoginUserHandler(c *gin.Context) {
	var login users.LoginUserRequest
	if err := c.ShouldBindJSON(&login); err != nil {
		log.Error().Msg(errorlogs.BindingJsonError)
		c.JSON(http.StatusBadRequest, gin.H{"error": errorlogs.BindingJsonError})
		return
	}

	dB := c.MustGet(constants.ConstantDb).(*gorm.DB)

	tx := dB.Begin()
	if tx.Error != nil {
		log.Error().Msgf(errorlogs.BeginSQLTransactionError, tx.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorlogs.BeginSQLTransactionError})
		return
	}

	defer tx.Rollback()

	user, err := users.GetUserByPhoneNumber(tx, login.PhoneNumber)
	if err != nil {
		tx.Rollback()
		log.Error().Msgf(errorlogs.GetUserByPhoneNumberError, login.PhoneNumber, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorlogs.GetUserByPhoneNumberError})
		return
	}

	if err = utils.VerifyPassword(user.PasswordHash, login.Password); err != nil {
		tx.Rollback()
		log.Error().Msg("Invalid password")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		message := fmt.Sprintf(errorlogs.CommitSQLTransactionError, constants.UserTable, err)
		log.Error().Msg("Error while commiting transaction to create user. " + message)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": message})
		return
	}

	token, ok := c.MustGet(constants.TokenMaker).(token.Maker)
	if !ok {
		log.Error().Msg("failed to retrieve token maker from context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	accessToken, err := token.CreateToken(
		login.PhoneNumber,
		user.ID,
		utils.GetAccessTokenDuration(),
	)
	if err != nil {
		log.Error().Msgf(errorlogs.TokenError, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	resp := users.LoginUserResponse{
		AccessToken: accessToken,
	}

	c.JSON(http.StatusOK, gin.H{"message": "login successful", "response": resp})
}
