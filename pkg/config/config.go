package config

import (
	"fmt"
	"spam-search/pkg/constants"
	"spam-search/pkg/constants/errorlogs"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var config *viper.Viper

func Init(env string) {
	var err error

	config = viper.New()
	config.SetConfigType(constants.DefaultConfigType)
	config.SetConfigName(env)
	config.AddConfigPath(constants.DefaultConfigPath)

	err = config.ReadInConfig()
	if err != nil {
		log.Fatal().Msg(fmt.Sprintf(errorlogs.ParsingError, err.Error()))
	}

}

func DBconfig() DB {
	return DB{
		URL:      config.GetString("db.url"),
		Username: config.GetString("db.username"),
		Password: config.GetString("db.password"),
		Database: config.GetString("db.name"),
	}
}

func GetConfig() *viper.Viper {
	return config
}

func GetSymmetricKey() string {
	return config.GetString("token.symmetric")
}

func GetAccessTokenDuration() string {
	return config.GetString("token.access.duration")
}
