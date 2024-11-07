package server

import (
	"net/http"
	"spam-search/pkg/constants"
	userController "spam-search/pkg/controller/users"
	"spam-search/pkg/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewRounter(dbConnection *gorm.DB) *gin.Engine {

	router := gin.New()

	router.Use(func(c *gin.Context) {
		c.Set(constants.ConstantDb, dbConnection)
	})

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	//TODO : apply logger middleware

	corsConfig := cors.DefaultConfig()
	router.Use(middleware.LoggerMiddleware())

	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowHeaders = []string{constants.Origin, constants.ContentLength, constants.ContentType, constants.Authorization}
	corsConfig.AllowMethods = []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}

	router.Use(cors.New(corsConfig))

	v1 := router.Group("/v1")
	{
		userGroup := v1.Group("/users")
		{
			userController := new(userController.UserController)
			userGroup.POST("/create", userController.CreateUserHandler)
			// userGroup.POST("/login", userController.LoginUserHandler)
		}
	}

	return router
}
