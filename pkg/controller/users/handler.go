package userController

import (
	"fmt"
	"net/http"
	"net/mail"
	"spam-search/pkg/constants"
	errorlogs "spam-search/pkg/constants/errorlogs"
	"spam-search/pkg/users"
	"spam-search/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func (con *UserController) CreateUserHandler(c *gin.Context) {
	var user users.User

	// Bind the incoming JSON data to the user struct
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Error().Msgf(errorlogs.BindingJsonError, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": errorlogs.BindingJsonError})
		return
	}

	// Validate the email format if provided
	if user.Email != "" {
		if _, err := mail.ParseAddress(user.Email); err != nil {
			log.Error().Msg(errorlogs.InvalidEmailFormatError)
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid email format"})
			return
		}
	}

	// Get the DB connection
	db := c.MustGet(constants.ConstantDb).(*gorm.DB)
	if db.Error != nil {
		message := fmt.Sprintf(errorlogs.GetDBError, db.Error)
		log.Error().Msg("Error while getting the db connection" + message)
		c.AbortWithStatusJSON(http.StatusInternalServerError, message)
		return
	}

	// Start a transaction to ensure atomicity
	tx := db.Begin()
	if tx.Error != nil {
		message := fmt.Sprintf(errorlogs.BeginSQLTransactionError, tx.Error)
		log.Error().Msg("Error while starting a transaction to insert user." + message)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": message})
		return
	}

	// Defer rollback in case of errors
	defer tx.Rollback()

	// Check if the user already exists by email or phone number
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

	// Hash the password before storing it
	hashedPassword, err := utils.HashPasswordArgon2(user.PasswordHash)
	if err != nil {
		tx.Rollback()
		log.Error().Msgf(errorlogs.HashedPasswordError, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": errorlogs.HashedPasswordError})
		return
	}

	// Set the hashed password to user
	user.PasswordHash = hashedPassword

	// Insert the user into the database
	if err := user.CreateUser(tx); err != nil {
		tx.Rollback()
		log.Error().Msg(errorlogs.InsertUserError + err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		message := fmt.Sprintf(errorlogs.CommitSQLTransactionError, constants.UserTable, err)
		log.Error().Msg("Error while commiting transaction to create user. " + message)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": message})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}
