package app

import (
	"syscall"
	"time"
	"context"
	"log"
	"os"
	"os/signal"
	"github.com/nullsec45/golang-news-api/lib/auth"
	"github.com/nullsec45/golang-news-api/config"
	"github.com/nullsec45/golang-news-api/lib/pagination"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/logger"
)


func RunServer(){
	cfg := config.NewConfig()
	_, err := cfg.ConnectionPostgres()

	if err != nil {
		log.Fatal("Error connection to database: %v", err)
		return
	}

	// Cloudflare R2
	cdfR2 := cfg.LoadAWSConfig()
	_ = s3.NewFromConfig(cdfR2)

	_ = auth.NewJwt(cfg)
	_ = pagination.NewPagination()

	app := fiber.New()
	app.Use(cors.New())
	app.Use(recover.New())
	// app.Use(logger)
	app.Use(logger.New(logger.Config{
		Format: "[${time}] %{ip} %{status} - %{latency} %{method} %{path}\n",
	}))

	_ = app.Group("/api")

	go func(){
		if cfg.App.AppPort == "" {
			cfg.App.AppPort = os.Getenv("APP_PORT")
		}

		err := app.Listen(":" + cfg.App.AppPort)
		
		if err != nil {
			log.Fatal("Error starting server:%v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)

	<-quit

	log.Println("server shutdown of 5 seconds")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	app.ShutdownWithContext(ctx)
}