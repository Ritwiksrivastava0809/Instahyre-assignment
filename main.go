package main

import (
	"flag"
	"fmt"
	"os"
	"spam-search/pkg/config"
	"spam-search/pkg/db"
	"spam-search/pkg/logger"
	"spam-search/pkg/server"

	"github.com/rs/zerolog/log"
)

func main() {

	logger.InitLogger()

	environment := flag.String("e", "development", "")

	flag.Usage = func() {
		fmt.Println("Usage: server -e {mode}")
		os.Exit(1)
	}

	config.Init(*environment)
	if *environment == "production" || *environment == "development" {
		log.Info().Msg("Environment is set to " + *environment)
	}

	dbConnection, err := db.NewSQLDB()
	if err != nil {
		log.Fatal().Msg("Can't Initialize DB " + err.Error())
		panic("Can't Initialize DB " + err.Error())
	}
	var developmentMode bool
	if *environment == "development" || *environment == "local" {
		developmentMode = true
	}
	dbConnection.Migrate(developmentMode)
	server.Init(dbConnection.GetDatabse())
}
