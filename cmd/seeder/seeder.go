package main

import (
	"log"

	"github.com/nullsec45/golang-news-api/config"
	dbseeds "github.com/nullsec45/golang-news-api/database/seeds"
)

func main() {
	// var cfgFile string


	// if cfgFile != "" {
	// 	viper.SetConfigFile(cfgFile)
	// }else{
	// 	viper.SetConfigFile(".env")
	// }

	config.Init() 
	cfg := config.NewConfig()

	pg, err := cfg.ConnectionPostgres()
	if err != nil {
		log.Fatal(err)
	}

	// KIRIMKAN *gorm.DB, bukan *config.Postgres
	dbseeds.SeedRoles(pg.DB)

	sqlDB, err := pg.DB.DB()
	if err == nil {
		defer sqlDB.Close()
	}
}
