package utils

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type QueryHelper struct {
	conn *pgx.Conn
}

func NewQueryHelper(conn *pgx.Conn) *QueryHelper {
	return &QueryHelper{conn: conn}
}

func (helper *QueryHelper) GetAllTables() ([]string, error) {
	var tables []string

	rows, err := helper.conn.Query(
		context.Background(),
		`SELECT
		 	table_schema || '.' || table_name 
		 FROM 
		 	information_schema.tables 
		 WHERE
		 	table_type = 'BASE TABLE' 
		 AND
		 	table_schema NOT IN ('pg_catalog', 'information_schema')`,
	)
	if err != nil {
		return tables, err
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		if err != nil {
			continue
		}
		tables = append(tables, tableName)
	}
	if err := rows.Err(); err != nil {
		return tables, err
	}

	return tables, nil
}
