package order

//go:generate mockgen -destination=./mocks/mock_OrderServicer.go -package=mocks github.com/enrico5b1b4/order-service/order OrderServicer

import (
	"database/sql"

	"github.com/enrico5b1b4/order-service/errors"
	"github.com/enrico5b1b4/order-service/orderprocess"
	uuid "github.com/satori/go.uuid"
)

type OrderServicer interface {
	GetOrderByID(uuid.UUID) (*Order, *errors.Error)
	GetOrders(string) ([]*Order, *errors.Error)
	CreateOrder(*Order) (*Order, *errors.Error)
	CompleteOrder(uuid.UUID, ProcessOrderStatus) *errors.Error
	ProcessOrder(ID uuid.UUID) *errors.Error
}

type OrderService struct {
	OrderStore          OrderStorer
	OrderProcessService orderprocess.OrderProcessor
}

func (s *OrderService) GetOrderByID(ID uuid.UUID) (*Order, *errors.Error) {
	o, err := s.OrderStore.GetOrderByOrderID(ID)
	if err != nil {
		return nil, mapStoreError(err)
	}

	return mapFromOrderDB(o), nil
}

func (s *OrderService) GetOrders(status string) ([]*Order, *errors.Error) {
	orderStatus, errValidateStatus := ValidateStatus(status)
	if errValidateStatus != nil {
		return nil, errValidateStatus
	}

	dbOrders, err := s.OrderStore.GetOrders(orderStatus)
	if err != nil {
		return nil, mapStoreError(err)
	}

	return mapFromOrderDBCollection(dbOrders), nil
}

func (s *OrderService) CreateOrder(o *Order) (*Order, *errors.Error) {
	o.Status = CREATED
	orderID, err := s.OrderStore.CreateOrder(mapToOrderDB(o))
	if err != nil {
		return nil, mapStoreError(err)
	}

	newOrder, err := s.OrderStore.GetOrderByID(orderID)
	if err != nil {
		return nil, mapStoreError(err)
	}

	return mapFromOrderDB(newOrder), nil
}

func (s *OrderService) ProcessOrder(ID uuid.UUID) *errors.Error {
	o, err := s.OrderStore.GetOrderByOrderID(ID)
	if err != nil {
		return mapStoreError(err)
	}

	if o.Status == PROCESSING {
		return OrderAlreadyBeingProcessedError
	}

	if o.Status == FULFILLED || o.Status == FAILED {
		return OrderAlreadyCompleteError
	}

	createOrderErr := s.OrderProcessService.CreateOrder(ID)
	if createOrderErr != nil {
		return errors.GeneralError
	}

	updateOrder := &OrderDB{
		ID:      o.ID,
		OrderID: o.OrderID,
		Status:  PROCESSING,
	}
	_, updateErr := s.OrderStore.UpdateOrder(updateOrder)
	if updateErr != nil {
		return mapStoreError(updateErr)
	}

	return nil
}

func (s *OrderService) CompleteOrder(ID uuid.UUID, status ProcessOrderStatus) *errors.Error {
	o, err := s.OrderStore.GetOrderByOrderID(ID)
	if err != nil {
		return mapStoreError(err)
	}

	orderStatus, mapError := mapOrderStatus(status)
	if mapError != nil {
		return mapError
	}

	o.Status = orderStatus

	_, updateErr := s.OrderStore.UpdateOrder(o)
	if updateErr != nil {
		return mapStoreError(updateErr)
	}

	return nil
}

func mapStoreError(err error) *errors.Error {
	switch err {
	case sql.ErrNoRows:
		return OrderNotFoundError
	case DBOrderAlreadyExistsError:
		return OrderAlreadyExistsError
	}

	return errors.GeneralError
}

func ValidateStatus(status string) (string, *errors.Error) {
	if status == "" ||
		status == CREATED ||
		status == PROCESSING ||
		status == FULFILLED ||
		status == FAILED {
		return status, nil
	}

	return "", InvalidOrdersStatusError
}
