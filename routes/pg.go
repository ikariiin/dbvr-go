package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ikariiin/dbvr-go/middleware"
	"github.com/ikariiin/dbvr-go/utils"
	"gorm.io/gorm"
)

type PgRoutes struct {
	db     *gorm.DB
	router *gin.Engine
}

func NewPgRoutes(db *gorm.DB, router *gin.Engine) *PgRoutes {
	return &PgRoutes{db: db, router: router}
}

func (r *PgRoutes) run(ctx *gin.Context) {
	user, err := utils.CurrentUser(ctx, r.db)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusOK, gin.H{"username": user.Username})
}

func (r *PgRoutes) RegisterRoutes() {
	group := r.router.Group("pg")
	group.Use(middleware.JwtAuthMiddleware())

	group.GET("run", r.run)
}
