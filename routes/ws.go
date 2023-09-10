package routes

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/ikariiin/dbvr-go/models"
	"github.com/ikariiin/dbvr-go/utils"
	"github.com/ikariiin/dbvr-go/wshandler"
	"github.com/jackc/pgx/v5"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WsRoutes struct {
	db     *gorm.DB
	router *gin.Engine
}

func NewWsRoutes(db *gorm.DB, router *gin.Engine) *WsRoutes {
	return &WsRoutes{db: db, router: router}
}

func (r *WsRoutes) registerWebSocket(ctx *gin.Context) {
	// user, err := utils.CurrentUser(ctx, r.db)
	token := ctx.Query("token")
	err := utils.ValidateRawToken(token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	user, err := utils.GetUserFromToken(token, r.db)

	connId, convertErr := strconv.Atoi(ctx.Param("connId"))
	if err != nil || convertErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	connections, err := models.GetUserConnections(r.db, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	findResult := slices.IndexFunc(connections, func(c models.Connection) bool {
		return c.ID == uint(connId)
	})
	if findResult == -1 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access to resource"})
		return
	}
	connectionStr := connections[findResult]

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	userDBConn, err := pgx.Connect(context.Background(), connectionStr.ConnString)
	if err != nil {
		ctx.AbortWithError(http.StatusNotFound, err)
	}
	defer userDBConn.Close(context.Background())

	for {
		var wsRequest wshandler.WsRequest
		err := conn.ReadJSON(&wsRequest)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			break
		}

		wshandler.HandleIncomingRequest(wsRequest, conn, userDBConn)
	}
}

func (r *WsRoutes) RegisterWsRoutes() {
	group := r.router.Group("ws")

	group.GET("/connect/:connId", r.registerWebSocket)
}
