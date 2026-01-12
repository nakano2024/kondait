package main

import (
	"log"

	"github.com/labstack/echo/v4"

	"kondait-backend/application/usecase"
	"kondait-backend/infra/config"
	"kondait-backend/infra/db"
	"kondait-backend/infra/repository"
	infrautil "kondait-backend/infra/util"
	"kondait-backend/web/handler"
	"kondait-backend/web/middleware"
)

func main() {
	cfgLoader := config.NewConfigLoader()
	cfg, err := cfgLoader.Load()
	if err != nil {
		log.Fatalln(err)
	}

	dbMigrator := db.NewDbMigrator()
	if err := dbMigrator.Migrate(cfg); err != nil {
		log.Fatalln(err)
	}

	dbInitializer := db.NewDbInitializer()
	db, err := dbInitializer.Open(cfg)
	if err != nil {
		log.Fatalln(err)
	}

	e := echo.New()
	healthHandler := handler.NewGetHealthHandler()
	e.GET("/health", healthHandler.Handle)

	var getPrincipalUsecase usecase.IGetPrincipalUsecase
	if cfg.Env == config.EnvDevelopment {
		getPrincipalUsecase = usecase.NewGetPrincipalUsecase(infrautil.NewPrincipalFetcherMock())
	} else {
		getPrincipalUsecase = usecase.NewGetPrincipalUsecase(infrautil.NewPrincipalFetcher())
	}

	getRecommendedCookingItemsUsecase := usecase.NewGetRecommendedCookingItemsUsecase(repository.NewRecommendedCookingItemRepository(db))
	getRecommendedCookingItemsHandler := handler.NewGetRecommendedCookingItemsHandler(getRecommendedCookingItemsUsecase)
	authApiGroup := e.Group("/api/private", middleware.AuthMiddleware(getPrincipalUsecase))
	authApiGroup.GET("/cooking-items/recommends", getRecommendedCookingItemsHandler.Handle)

	if err := e.Start(":" + cfg.Port); err != nil {
		log.Fatalln(err)
	}
}
