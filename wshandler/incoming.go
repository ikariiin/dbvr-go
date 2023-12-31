package wshandler

import (
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5"
)

type RequestType string

const (
	Query   RequestType = "query"
	InBuilt RequestType = "inBuilt"
)

type InBuiltType string

const (
	UserTables InBuiltType = "userTables"
)

type WsRequest struct {
	Name         string      `json:"name"`
	ID           string      `json:"requestId"`
	RequestType  RequestType `json:"requestType"`
	RequestValue string      `json:"requestValue"`
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
