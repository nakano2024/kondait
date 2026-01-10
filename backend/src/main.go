package main

import (
	"log"

	"github.com/labstack/echo/v4"

	"kondait-backend/handler"
	"kondait-backend/infra/config"
	"kondait-backend/infra/db"
)

func main() {
	cfgLoader := config.NewConfigLoader()
	cfg, err := cfgLoader.Load()
	if err != nil {
		log.Fatalln(err)
	}

	_, err = db.DbInitializer{}.Open(cfg)
	if err != nil {
		log.Fatalln(err)
	}

	e := echo.New()
	healthHandler := handler.NewGetHealthHandler()
	e.GET("/health", healthHandler.Handle)

	if err := e.Start(":" + cfg.Port); err != nil {
		log.Fatalln(err)
	}
}
