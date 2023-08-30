package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ikariiin/dbvr-go/middleware"
	"github.com/ikariiin/dbvr-go/models"
	"github.com/ikariiin/dbvr-go/routes"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	db, err := models.SetupDB()
	if err != nil {
		log.Println("Could not initialize DB")
	}
	return db
}

func InitRouter() *gin.Engine {
	r := gin.Default()
	db := InitDB()

	authRoutes := routes.NewAuthRoutes(db, r)
	authRoutes.RegisterAuthRoutes()

	authGroup := r.Group("/authorized")
	authGroup.Use(middleware.JwtAuthMiddleware())
	authGroup.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"allowed": true})
	})

	return r
}

func main() {
	r := InitRouter()

	r.Run()
}
