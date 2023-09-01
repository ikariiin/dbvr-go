package wshandler

import (
	"github.com/ikariiin/dbvr-go/utils"
	"github.com/jackc/pgx/v5"
)

type WsResponse struct {
	Name       string      `json:"name"`
	ResponseTo *string     `json:"response-to"`
	Value      interface{} `json:"value"`
}

func GetUserTables(db *pgx.Conn) ([]string, error) {
	h := utils.NewQueryHelper(db)
	return h.GetAllTables()
}
