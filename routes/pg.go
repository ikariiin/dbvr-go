package routes

import "github.com/gin-gonic/gin"

func RegisterPGRoutes(r *gin.Engine) {
	pgRoutes := r.Group("pg")
	{
		pgRoutes.GET("/test", func(g *gin.Context) {
			g.JSON(200, gin.H{
				"test": "toast",
			})
		})
	}
}
