package server

import (
	"fmt"
	"net/http"
	"spam-search/pkg/config"
	"spam-search/pkg/constants"
	errorlogs "spam-search/pkg/constants/errorlogs"
	spamReportsController "spam-search/pkg/controller/spamReports"
	userController "spam-search/pkg/controller/users"
	"spam-search/pkg/middleware"
	"spam-search/pkg/token"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Token struct {
	tokenMaker token.Maker
}

func NewRounter(dbConnection *gorm.DB) (*gin.Engine, error) {
	tokenMaker, err := token.NewJWTMAKER(config.GetSymmetricKey())
	if err != nil {
		return nil, fmt.Errorf(errorlogs.TokenError, err)
	}

	router := gin.New()

	router.Use(func(c *gin.Context) {
		c.Set(constants.ConstantDb, dbConnection)
		c.Set(constants.TokenMaker, tokenMaker)
	})

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	corsConfig := cors.DefaultConfig()
	router.Use(middleware.LoggerMiddleware())

	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowHeaders = []string{constants.Origin, constants.ContentLength, constants.ContentType, constants.Authorization}
	corsConfig.AllowMethods = []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}

	router.Use(cors.New(corsConfig))

	authMiddleWare := middleware.AuthTokenMiddleware(tokenMaker)

	v1 := router.Group("/v1")
	{
		userGroup := v1.Group("/users")
		{
			userController := new(userController.UserController)
			userGroup.POST("/create", userController.CreateUserHandler)
			userGroup.POST("/login", userController.LoginUserHandler)
		}

		spamGroup := v1.Group("/spam")
		spamGroup.Use(authMiddleWare)
		{
			spamController := new(spamReportsController.SpamReportsController)
			spamGroup.POST("/report", spamController.ReportSpamHandler)
			spamGroup.GET("/search/name", spamController.SearchNameHandler)
			spamGroup.GET("/search/phone", spamController.SearchPhoneHandler)
		}
	}

	return router, nil
}
