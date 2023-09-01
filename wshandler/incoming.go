package wshandler

import (
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5"
)

type RequestType string

const (
	Query   RequestType = "query"
	InBuilt RequestType = "in-built"
)

type InBuiltType string

const (
	UserTables InBuiltType = "user-tables"
)

type WsRequest struct {
	Name         string      `json:"name"`
	ID           string      `json:"request-id"`
	RequestType  RequestType `json:"request-type"`
	RequestValue string      `json:"request-value"`
}

func handleInBuiltRequests(request WsRequest, db *pgx.Conn) (WsResponse, error) {
	response := WsResponse{ResponseTo: &request.ID, Name: "InBuilt Response"}

	switch request.RequestValue {
	case string(UserTables):
		tables, err := GetUserTables(db)
		if err != nil {
			// TODO: Error handling
			break
		}
		response.Value = tables
	}

	return response, nil
}

func HandleIncomingRequest(request WsRequest, conn *websocket.Conn, db *pgx.Conn) {
	switch request.RequestType {
	case Query:
		// TODO:
		break
	case InBuilt:
		// TODO:
		response, err := handleInBuiltRequests(request, db)
		if err != nil {
			// TODO: Error handling
		}

		conn.WriteJSON(response)
		break
	}
}
