package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WsRoutes struct {
	db     *gorm.DB
	router *gin.Engine
}

func NewWsRoutes(db *gorm.DB, router *gin.Engine) *WsRoutes {
	return &WsRoutes{db: db, router: router}
}
