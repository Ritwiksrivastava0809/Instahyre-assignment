package server

import (
	"fmt"
	"spam-search/pkg/config"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func Init(dbConnection *gorm.DB) {
	config := config.GetConfig()
	router, err := NewRounter(dbConnection)
	if err != nil {
		log.Fatal().Msg(fmt.Sprintf("error while initializing the server :: %s", err))
	}
	log.Debug().Msg(config.GetString("server.port"))
	serverUrl := fmt.Sprintf("%s:%s", config.GetString("server.host"), config.GetString("server.port"))
	router.Run(serverUrl)
}
