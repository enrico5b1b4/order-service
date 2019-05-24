package order

//go:generate mockgen -destination=./mocks/mock_OrderStorer.go -package=mocks github.com/enrico5b1b4/order-service/order OrderStorer

import (
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
)

type OrderDB struct {
	ID        int       `db:"id"`
	OrderID   uuid.UUID `db:"order_id" sql:",type:uuid"`
	Status    string    `db:"order_status"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type OrderStorer interface {
	GetOrderByID(int) (*OrderDB, error)
	GetOrderByOrderID(uuid.UUID) (*OrderDB, error)
	GetOrders(string) ([]*OrderDB, error)
	CreateOrder(*OrderDB) (int, error)
	UpdateOrder(*OrderDB) (int, error)
}

var DBOrderAlreadyExistsError = errors.New("order already exists")

type OrderStore struct {
	DB *sqlx.DB
}

func (s *OrderStore) GetOrderByID(ID int) (*OrderDB, error) {
	var order OrderDB
	err := s.DB.Get(&order, `
		SELECT id, order_id, order_status
		FROM orders 
		WHERE id = $1`, ID,
	)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (s *OrderStore) GetOrderByOrderID(orderID uuid.UUID) (*OrderDB, error) {
	var order OrderDB
	err := s.DB.Get(&order, `
		SELECT id, order_id, order_status, created_at, updated_at
		FROM orders 
		WHERE order_id = $1`, orderID.String(),
	)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (s *OrderStore) GetOrders(filter string) ([]*OrderDB, error) {
	var orders []*OrderDB
	var err error

	if filter != "" {
		err = s.DB.Select(&orders, `
		SELECT id, order_id, order_status, created_at, updated_at
		FROM orders
		WHERE order_status = $1`, filter)
	} else {
		err = s.DB.Select(&orders,
			`
		SELECT 
			id, order_id, order_status, created_at, updated_at
		FROM 
			orders`)
	}
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *OrderStore) CreateOrder(o *OrderDB) (int, error) {
	tx, err := s.DB.Beginx()
	if err != nil {
		return 0, err
	}
	o.CreatedAt = time.Now()
	o.UpdatedAt = time.Now()

	stmt, err := tx.Preparex(`	
		INSERT INTO orders 
			(order_id, order_status, created_at, updated_at)
		VALUES 
			($1, $2, $3, $4) 
		RETURNING id`)
	if err != nil {
		errRB := tx.Rollback()
		if errRB != nil {
			return 0, errRB
		}
		return 0, err
	}
	defer stmt.Close()

	var ID int
	err = stmt.Get(&ID, o.OrderID, o.Status, o.CreatedAt, o.UpdatedAt)
	if err != nil {
		errRB := tx.Rollback()
		if errRB != nil {
			return 0, errRB
		}
		return 0, mapCreateOrderError(err)
	}

	errC := tx.Commit()
	if errC != nil {
		return 0, errC
	}

	return ID, nil
}

func (s *OrderStore) UpdateOrder(o *OrderDB) (int, error) {
	tx, err := s.DB.Beginx()
	if err != nil {
		return 0, err
	}
	o.UpdatedAt = time.Now()

	stmt, err := tx.Preparex(`	
		UPDATE orders 
		SET order_status = $1, updated_at = $2
		WHERE id = $3
		RETURNING id`)
	if err != nil {
		errRB := tx.Rollback()
		if errRB != nil {
			return 0, errRB
		}
		return 0, err
	}
	defer stmt.Close()

	var ID int
	err = stmt.Get(&ID, o.Status, o.UpdatedAt, o.ID)
	if err != nil {
		errRB := tx.Rollback()
		if errRB != nil {
			return 0, errRB
		}
		return 0, err
	}

	errC := tx.Commit()
	if errC != nil {
		return 0, errC
	}

	return ID, nil
}

func mapCreateOrderError(err error) error {
	if pgErr, ok := err.(*pq.Error); ok {
		// 23505 - unique_violation
		if pgErr.Code == "23505" {
			return DBOrderAlreadyExistsError
		}
	}

	return err
}
