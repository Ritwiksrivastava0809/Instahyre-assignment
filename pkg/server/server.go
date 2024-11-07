package server

import (
	"fmt"
	"spam-search/pkg/config"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func Init(dbConnection *gorm.DB) {
	config := config.GetConfig()
	router := NewRounter(dbConnection)
	log.Debug().Msg(config.GetString("server.port"))
	serverUrl := fmt.Sprintf("%s:%s", config.GetString("server.host"), config.GetString("server.port"))
	router.Run(serverUrl)
}
