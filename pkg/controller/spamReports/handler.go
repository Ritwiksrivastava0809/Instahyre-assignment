package spamReportsController

import (
	"fmt"
	"net/http"
	"spam-search/pkg/constants"
	"spam-search/pkg/constants/errorlogs"
	spamreports "spam-search/pkg/spamReports"
	"spam-search/pkg/token"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func (con *SpamReportsController) ReportSpamHandler(c *gin.Context) {
	var reportRequest spamreports.SpamReportRequest

	if err := c.ShouldBindJSON(&reportRequest); err != nil {
		log.Error().Msgf(errorlogs.BindingJsonError, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf(errorlogs.BindingJsonError, err.Error())})
		return
	}

	db := c.MustGet(constants.ConstantDb).(*gorm.DB)

	spamReport, err := spamreports.ReportSpam(db, reportRequest.PhoneNumber, reportRequest.Name)
	if err != nil {
		log.Error().Msgf(errorlogs.ReportSpamError, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf(errorlogs.ReportSpamError, err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Spam reported successfully", constants.SpamReport: spamReport})
}

func (con *SpamReportsController) SearchNameHandler(c *gin.Context) {
	query := c.DefaultQuery(constants.Name, "")
	if query == "" {
		log.Error().Msg(errorlogs.MissingQueryParameterError)
		c.JSON(http.StatusBadRequest, gin.H{"message": errorlogs.MissingQueryParameterError})
		return
	}

	db := c.MustGet(constants.ConstantDb).(*gorm.DB)
	if db.Error != nil {
		message := fmt.Sprintf(errorlogs.GetDBError, db.Error)
		log.Error().Msg(message)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": message})
		return
	}

	results, err := spamreports.GetSpamReportsByName(db, query)
	if err != nil {
		log.Error().Msgf(errorlogs.DatabaseQueryError, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error searching for name"})
		return
	}

	if len(results) == 0 {
		results, err = spamreports.GetSpamReportsByName(db, "%"+query+"%")
		if err != nil {
			log.Error().Msgf(errorlogs.DatabaseQueryError, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error searching for name"})
			return
		}
	}

	var response []gin.H
	for _, result := range results {
		email := ""
		authPayload := c.MustGet(constants.AuthorizationPayloadKey).(*token.Payload)
		currentUserID := authPayload.UserID

		email, err = spamreports.GetUserContactEmail(db, currentUserID, result.PhoneNumber)
		if err != nil {
			log.Error().Msgf("Error retrieving user contact email: %v", err)
		}

		response = append(response, gin.H{
			constants.Name:           result.Name,
			constants.PhoneNumber:    result.PhoneNumber,
			constants.SpamLikelihood: result.SpamLikelihood,
			constants.Email:          email,
		})
	}

	c.JSON(http.StatusOK, gin.H{"results": response})
}

func (con *SpamReportsController) SearchPhoneHandler(c *gin.Context) {
	phoneNumber := c.DefaultQuery(constants.PhoneNumber, "")
	if phoneNumber == "" {
		log.Error().Msg(errorlogs.MissingQueryParameterError)
		c.JSON(http.StatusBadRequest, gin.H{"message": errorlogs.MissingQueryParameterError})
		return
	}

	db := c.MustGet(constants.ConstantDb).(*gorm.DB)
	if db.Error != nil {
		message := fmt.Sprintf(errorlogs.GetDBError, db.Error)
		log.Error().Msg(message)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": message})
		return
	}

	user, err := spamreports.GetUserByPhoneNumber(db, phoneNumber)
	if err == nil {
		var result spamreports.GlobalSpam
		err := db.Where("phone_number = ?", phoneNumber).First(&result).Error
		if err != nil {
			log.Error().Msgf(errorlogs.DatabaseQueryError, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving global spam record"})
			return
		}

		email := user.Email
		c.JSON(http.StatusOK, gin.H{
			constants.Name:           user.Name,
			constants.PhoneNumber:    user.PhoneNumber,
			constants.SpamLikelihood: result.SpamLikelihood,
			constants.Email:          email,
		})
		return
	}

	results, err := spamreports.GetSpamReportsByPhoneNumber(db, phoneNumber)
	if err != nil {
		log.Error().Msgf(errorlogs.DatabaseQueryError, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error searching for phone number"})
		return
	}

	if len(results) == 0 {
		results, err = spamreports.GetSpamReportsByPhoneNumber(db, "%"+phoneNumber+"%")
		if err != nil {
			log.Error().Msgf(errorlogs.DatabaseQueryError, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error searching for phone number"})
			return
		}
	}

	var response []gin.H
	for _, result := range results {
		email := ""

		authPayload := c.MustGet(constants.AuthorizationPayloadKey).(*token.Payload)
		currentUserID := authPayload.UserID

		email, err = spamreports.GetUserContactEmail(db, currentUserID, result.PhoneNumber)
		if err != nil {
			log.Error().Msgf("Error retrieving user contact email: %v", err)
		}

		response = append(response, gin.H{
			constants.Name:           result.Name,
			constants.PhoneNumber:    result.PhoneNumber,
			constants.SpamLikelihood: result.SpamLikelihood,
			constants.Email:          email,
		})
	}

	c.JSON(http.StatusOK, gin.H{"results": response})
}
