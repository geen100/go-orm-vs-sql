package main

import (
	"go-orm-vs-sql/config"
	"go-orm-vs-sql/database"
	"go-orm-vs-sql/handlers"
	"go-orm-vs-sql/middleware"
	"go-orm-vs-sql/repositories"
	"log"

	"github.com/gin-gonic/gin"
)

func setupRoutes(userRepo repositories.UserRepository) *gin.Engine {
	r := gin.Default()

	// CORS設定
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	authHandler := handlers.NewAuthHandler(userRepo)

	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		api.GET("/me", authHandler.Me)
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return r
}

func main() {
	config.Init()

	gormDB, err := database.Init()
	if err != nil {
		log.Fatal("Failed to initialize GORM database:", err)
	}

	userRepo := repositories.NewGORMUserRepository(gormDB)
	log.Println("Using GORM repository")

	r := setupRoutes(userRepo)

	port := "8080"
	log.Printf("Server starting on port %s with GORM repository", port)
	r.Run(":" + port)
}
