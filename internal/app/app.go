package app

import (
	"github.com/nullsec45/golang-news-api/config"
	"github.com/rs/zerolog/log"
)


func RunServer(){
	cfg := config.NewConfig()
	_, err := cfg.ConnectionPostgres()

	if err != nil {
		log.Fatal().Msgf("Error connection to database: %v", err)
		return
	}
}