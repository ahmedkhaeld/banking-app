package main

import (
	"log"
	"os"

	"github.com/ahmedkhaeld/banking-app/db"
	_ "github.com/ahmedkhaeld/banking-app/docs" // Import the generated docs
	"github.com/ahmedkhaeld/banking-app/internal/user"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @contact.name               API Support
// @contact.url                http://www.swagger.io/support
// @contact.email              support@swagger.io
// @securityDefinitions.apikey JWT
// @in                         header
// @name                       Authorization
// @BasePath  /api/v1/

// @license.name Apache 2.0
// @license.url  http://www.apache.org/licenses/LICENSE-2.0.html
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

	server.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	server.Run(":" + os.Getenv("PORT"))
}
