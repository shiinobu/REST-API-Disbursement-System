package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"rest-api-disbursement-system/internal/config"
	"rest-api-disbursement-system/internal/handlers"
	"rest-api-disbursement-system/internal/middleware"
	"rest-api-disbursement-system/internal/models"
	"rest-api-disbursement-system/internal/repository"
	"rest-api-disbursement-system/internal/services"
)

func main() {
	cfg := config.Load()

	db, err := config.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.Disbursement{}); err != nil {
		log.Fatalf("database migration failed: %v", err)
	}

	if err := config.SeedUsers(db); err != nil {
		log.Fatalf("database seeding failed: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	disbursementRepo := repository.NewDisbursementRepository(db)

	authService := services.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTExpiresInHours)
	disbursementService := services.NewDisbursementService(disbursementRepo)

	authHandler := handlers.NewAuthHandler(authService)
	disbursementHandler := handlers.NewDisbursementHandler(disbursementService)

	router := gin.Default()

	api := router.Group("/api")
	authHandler.MountRoutes(api.Group("/auth"))

	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	disbursementHandler.MountRoutes(protected.Group("/disbursements"))

	log.Printf("server running on port %s", cfg.AppPort)
	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
