package order

import (
	"github.com/enrico5b1b4/order-service/errors"
	uuid "github.com/satori/go.uuid"
)

type OrderStatus = string

const (
	CREATED    OrderStatus = "CREATED"
	PROCESSING OrderStatus = "PROCESSING"
	FULFILLED  OrderStatus = "FULFILLED"
	FAILED     OrderStatus = "FAILED"
	UNKNOWN    OrderStatus = "UNKNOWN"
)

type ProcessOrderStatus = string

const (
	ORDER_SUCCEEDED ProcessOrderStatus = "SUCCEEDED"
	ORDER_FAILED    ProcessOrderStatus = "FAILED"
)

type Order struct {
	ID     uuid.UUID   `json:"id"`
	Status OrderStatus `json:"status"`
}

func mapFromOrderDB(o *OrderDB) *Order {
	return &Order{
		ID:     o.OrderID,
		Status: o.Status,
	}
}

func mapFromOrderDBCollection(dbOrders []*OrderDB) []*Order {
	orders := make([]*Order, len(dbOrders))
	for i := range dbOrders {
		orders[i] = mapFromOrderDB(dbOrders[i])
	}
	return orders
}

func mapToOrderDB(o *Order) *OrderDB {
	return &OrderDB{
		OrderID: o.ID,
		Status:  o.Status,
	}
}

func mapOrderStatus(status ProcessOrderStatus) (OrderStatus, *errors.Error) {
	if status == ORDER_SUCCEEDED {
		return FULFILLED, nil
	}

	if status == ORDER_FAILED {
		return FAILED, nil
	}

	return UNKNOWN, UnknownProcessOrderStatusError
}
