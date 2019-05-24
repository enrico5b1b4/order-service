package test

import (
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
)

func InsertOrder(testDB *sqlx.DB, orderID uuid.UUID, status string) int {
	var newOrderID int
	testDB.QueryRow(
		`INSERT INTO orders
					(order_id, order_status) 
				VALUES 
					($1, $2) 
				RETURNING id`, orderID.String(), status,
	).Scan(&newOrderID) // #nosec
	return newOrderID
}

func TruncateOrders(testDB *sqlx.DB) {
	testDB.MustExec(`TRUNCATE orders`)
}
