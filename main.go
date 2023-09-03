package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ikariiin/dbvr-go/models"
	"github.com/ikariiin/dbvr-go/routes"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	db, err := models.SetupDB()
	if err != nil {
		log.Println("Could not initialize DB")
	}
	return db
}

func CorsConfig() cors.Config {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"https://dbvr.ikariiin.xyz"}
	if os.Getenv("GIN_MODE") != "release" {
		config.AllowOrigins = append(config.AllowOrigins, "http://localhost:5173")
	}

	return config
}

func InitRouter() *gin.Engine {
	r := gin.Default()
	db := InitDB()

	r.Use(cors.New(CorsConfig()))

	authRoutes := routes.NewAuthRoutes(db, r)
	authRoutes.RegisterAuthRoutes()
	pgRoutes := routes.NewPgRoutes(db, r)
	pgRoutes.RegisterRoutes()
	wsRoutes := routes.NewWsRoutes(db, r)
	wsRoutes.RegisterWsRoutes()

	return r
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Could not load .env file")
	}

	r := InitRouter()

	r.Run()
}
