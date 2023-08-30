package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ikariiin/dbvr-go/middleware"
	"github.com/ikariiin/dbvr-go/models"
	"github.com/ikariiin/dbvr-go/utils"
	"gorm.io/gorm"
)

type CreateConnectionDTO struct {
	ConnectionString string `json:"connection-string" binding:"required"`
}

type PgRoutes struct {
	db     *gorm.DB
	router *gin.Engine
}

func NewPgRoutes(db *gorm.DB, router *gin.Engine) *PgRoutes {
	return &PgRoutes{db: db, router: router}
}

func (r *PgRoutes) getConnections(ctx *gin.Context) {
	user, err := utils.CurrentUser(ctx, r.db)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	connections, err := models.GetUserConnections(r.db, user)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch user connection strings"})
	}

	ctx.JSON(http.StatusOK, connections)
}

func (r *PgRoutes) createConnection(ctx *gin.Context) {
	user, err := utils.CurrentUser(ctx, r.db)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	var input CreateConnectionDTO
	if err := ctx.ShouldBind(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.CreateUserConnection(r.db, user, input.ConnectionString); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "message": "Created connection string for user"})

}

func (r *PgRoutes) RegisterRoutes() {
	group := r.router.Group("pg")
	group.Use(middleware.JwtAuthMiddleware())

	group.GET("connection", r.getConnections)
	group.POST("connection", r.createConnection)
}
