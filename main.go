package main

import (
	"log"
	"os"

	"github.com/ahmedkhaeld/banking-app/db"
	"github.com/ahmedkhaeld/banking-app/internal/user"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	server := gin.New()
	server.Use(gin.Recovery())

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AllowHeaders = []string{"Content-Type", "Authorization"}
	server.Use(cors.New(config))

	if os.Getenv("GIN_MODE") == "debug" {
		server.Use(gin.Logger())
		server.Use(gin.Recovery())
	}
	gin.SetMode(os.Getenv("GIN_MODE"))

	if err := db.Open(os.Getenv("DB_SOURCE")); err != nil {
		log.Fatal("Error opening database: ", err)
	}

	if err := db.RunMigrations(); err != nil {
		log.Fatal("Error running migrations: ", err)
	}

	//TODO: ADD SWAGGER

	server.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "OK"})
	})

	apiV1 := server.Group("/api/v1")

	userGroup := apiV1.Group("/users")
	user.RegisterRoutes(userGroup)

	server.Run(":" + os.Getenv("PORT"))
}
