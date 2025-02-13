package main

import (
	"log"
	_ "main/docs"
	"main/internal/config"
	"main/internal/repositories"
	"main/internal/routes"
	"main/internal/store"
	"main/pkg/middleware"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	cfg, err := config.NewFromEnv()
	if err != nil {
		log.Fatalf("Can`t read congif from ENV: %v", err)
	}

	newStore, err := store.NewStore(*cfg)
	if err != nil {
		log.Fatalf("Can`t create store: %v", err)
	}
	catRepo := repositories.NewCatRepository(*newStore)
	missionRepo := repositories.NewMissionRepository(*newStore)

	r := routes.SetupRouter(catRepo, missionRepo)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Use(middleware.Logger())

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
