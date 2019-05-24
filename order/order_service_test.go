package order_test

import (
	"database/sql"
	"errors"
	"testing"

	errs "github.com/enrico5b1b4/order-service/errors"
	"github.com/enrico5b1b4/order-service/order"
	orderMocks "github.com/enrico5b1b4/order-service/order/mocks"
	orderprocessMocks "github.com/enrico5b1b4/order-service/orderprocess/mocks"
	"github.com/golang/mock/gomock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

var OrderID = uuid.Must(uuid.FromString("f6c7d890-ca9b-4147-aa5b-2a41cb17f95a"))

func TestOrderService_GetOrderByID_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockStore := orderMocks.NewMockOrderStorer(mockCtrl)
	mockStore.
		EXPECT().
		GetOrderByOrderID(OrderID).
		Return(&order.OrderDB{ID: 1, OrderID: OrderID, Status: "CREATED"}, nil).
		Times(1)
	mockOrderProcessService := orderprocessMocks.NewMockOrderProcessor(mockCtrl)
	orderService := &order.OrderService{
		OrderStore:          mockStore,
		OrderProcessService: mockOrderProcessService,
	}

	o, err := orderService.GetOrderByID(OrderID)

	assert.Equal(t, &order.Order{ID: OrderID, Status: order.CREATED}, o)
	assert.Nil(t, err)
}

func TestOrderService_GetOrderByID_ErrorGeneralError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockStore := orderMocks.NewMockOrderStorer(mockCtrl)
	mockStore.
		EXPECT().
		GetOrderByOrderID(OrderID).
		Return(nil, errors.New("ERROR")).
		Times(1)
	mockOrderProcessService := orderprocessMocks.NewMockOrderProcessor(mockCtrl)
	orderService := &order.OrderService{
		OrderStore:          mockStore,
		OrderProcessService: mockOrderProcessService,
	}

	o, err := orderService.GetOrderByID(OrderID)

	assert.Equal(t, errs.GeneralError, err)
	assert.Nil(t, o)
}

func TestOrderService_GetOrderByID_ErrorOrderNotFound(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockStore := orderMocks.NewMockOrderStorer(mockCtrl)
	mockStore.
		EXPECT().
		GetOrderByOrderID(OrderID).
		Return(nil, sql.ErrNoRows).
		Times(1)
	mockOrderProcessService := orderprocessMocks.NewMockOrderProcessor(mockCtrl)
	orderService := &order.OrderService{
		OrderStore:          mockStore,
		OrderProcessService: mockOrderProcessService,
	}

	o, err := orderService.GetOrderByID(OrderID)

	assert.Equal(t, order.OrderNotFoundError, err)
	assert.Nil(t, o)
}

func TestOrderService_GetOrders_SuccessWithoutFilterStatus(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockStore := orderMocks.NewMockOrderStorer(mockCtrl)
	mockStore.
		EXPECT().
		GetOrders("").
		Return([]*order.OrderDB{{ID: 1, OrderID: OrderID, Status: "CREATED"}}, nil).
		Times(1)
	mockOrderProcessService := orderprocessMocks.NewMockOrderProcessor(mockCtrl)
	orderService := &order.OrderService{
		OrderStore:          mockStore,
		OrderProcessService: mockOrderProcessService,
	}

	o, err := orderService.GetOrders("")

	assert.Equal(t, []*order.Order{{ID: OrderID, Status: order.CREATED}}, o)
	assert.Nil(t, err)
}

func TestOrderService_GetOrders_SuccessWithFilterStatus(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockStore := orderMocks.NewMockOrderStorer(mockCtrl)
	mockStore.
		EXPECT().
		GetOrders("CREATED").
		Return([]*order.OrderDB{{ID: 1, OrderID: OrderID, Status: "CREATED"}}, nil).
		Times(1)
	mockOrderProcessService := orderprocessMocks.NewMockOrderProcessor(mockCtrl)
	orderService := &order.OrderService{
		OrderStore:          mockStore,
		OrderProcessService: mockOrderProcessService,
	}

	o, err := orderService.GetOrders("CREATED")

	assert.Equal(t, []*order.Order{{ID: OrderID, Status: order.CREATED}}, o)
	assert.Nil(t, err)
}

func TestOrderService_GetOrders_ErrorInvalidStatus(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockStore := orderMocks.NewMockOrderStorer(mockCtrl)
	mockOrderProcessService := orderprocessMocks.NewMockOrderProcessor(mockCtrl)
	orderService := &order.OrderService{
		OrderStore:          mockStore,
		OrderProcessService: mockOrderProcessService,
	}

	o, err := orderService.GetOrders("NOT_VALID")

	assert.Equal(t, order.InvalidOrdersStatusError, err)
	assert.Nil(t, o)
}

func TestOrderService_GetOrders_ErrorGeneralError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockStore := orderMocks.NewMockOrderStorer(mockCtrl)
	mockStore.
		EXPECT().
		GetOrders("").
		Return(nil, errors.New("ERROR")).
		Times(1)
	mockOrderProcessService := orderprocessMocks.NewMockOrderProcessor(mockCtrl)
	orderService := &order.OrderService{
		OrderStore:          mockStore,
		OrderProcessService: mockOrderProcessService,
	}

	o, err := orderService.GetOrders("")

	assert.Equal(t, errs.GeneralError, err)
	assert.Nil(t, o)
}

func TestOrderService_CreateOrder_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockStore := orderMocks.NewMockOrderStorer(mockCtrl)
	mockStore.
		EXPECT().
		CreateOrder(&order.OrderDB{
			OrderID: OrderID,
			Status:  "CREATED",
		}).
		Return(123, nil).
		Times(1)
	mockStore.
		EXPECT().
		GetOrderByID(123).
		Return(&order.OrderDB{
			ID:      123,
			OrderID: OrderID,
			Status:  "CREATED",
		}, nil).
		Times(1)
	mockOrderProcessService := orderprocessMocks.NewMockOrderProcessor(mockCtrl)
	orderService := &order.OrderService{
		OrderStore:          mockStore,
		OrderProcessService: mockOrderProcessService,
	}

	o, err := orderService.CreateOrder(&order.Order{ID: OrderID})

	assert.Equal(t, &order.Order{ID: OrderID, Status: order.CREATED}, o)
	assert.Nil(t, err)
}

func TestOrderService_CreateOrder_ErrorOrderAlreadyExistsError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockStore := orderMocks.NewMockOrderStorer(mockCtrl)
	mockStore.
		EXPECT().
		CreateOrder(&order.OrderDB{
			OrderID: OrderID,
			Status:  "CREATED",
		}).
		Return(0, order.DBOrderAlreadyExistsError).
		Times(1)
	mockOrderProcessService := orderprocessMocks.NewMockOrderProcessor(mockCtrl)
	orderService := &order.OrderService{
		OrderStore:          mockStore,
		OrderProcessService: mockOrderProcessService,
	}

	o, err := orderService.CreateOrder(&order.Order{ID: OrderID})

	assert.Equal(t, order.OrderAlreadyExistsError, err)
	assert.Nil(t, o)
}

func TestOrderService_CreateOrder_ErrorGeneralErrorCreatingOrder(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockStore := orderMocks.NewMockOrderStorer(mockCtrl)
	mockStore.
		EXPECT().
		CreateOrder(&order.OrderDB{
			OrderID: OrderID,
			Status:  "CREATED",
		}).
		Return(0, errors.New("ERROR")).
		Times(1)
	mockOrderProcessService := orderprocessMocks.NewMockOrderProcessor(mockCtrl)
	orderService := &order.OrderService{
		OrderStore:          mockStore,
		OrderProcessService: mockOrderProcessService,
	}

	o, err := orderService.CreateOrder(&order.Order{ID: OrderID})

	assert.Equal(t, errs.GeneralError, err)
	assert.Nil(t, o)
}

func TestOrderService_CreateOrder_ErrorGeneralErrorFetchingOrder(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockStore := orderMocks.NewMockOrderStorer(mockCtrl)
	mockStore.
		EXPECT().
		CreateOrder(&order.OrderDB{
			OrderID: OrderID,
			Status:  "CREATED",
		}).
		Return(123, nil).
		Times(1)
	mockStore.
		EXPECT().
		GetOrderByID(123).
		Return(nil, errors.New("ERROR")).
		Times(1)
	mockOrderProcessService := orderprocessMocks.NewMockOrderProcessor(mockCtrl)
	orderService := &order.OrderService{
		OrderStore:          mockStore,
		OrderProcessService: mockOrderProcessService,
	}

	o, err := orderService.CreateOrder(&order.Order{ID: OrderID})

	assert.Equal(t, errs.GeneralError, err)
	assert.Nil(t, o)
}

func TestOrderService_ProcessOrder_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockStore := orderMocks.NewMockOrderStorer(mockCtrl)
	mockStore.
		EXPECT().
		GetOrderByOrderID(OrderID).
		Return(&order.OrderDB{
			ID:      123,
			OrderID: OrderID,
			Status:  "CREATED",
		}, nil).
		Times(1)
	mockStore.
		EXPECT().
		UpdateOrder(&order.OrderDB{
			ID:      123,
			OrderID: OrderID,
			Status:  "PROCESSING",
		}).
		Return(123, nil).
		Times(1)
	mockOrderProcessService := orderprocessMocks.NewMockOrderProcessor(mockCtrl)
	mockOrderProcessService.
		EXPECT().
		CreateOrder(OrderID).
		Return(nil).
		Times(1)
	orderService := &order.OrderService{
		OrderStore:          mockStore,
		OrderProcessService: mockOrderProcessService,
	}

	err := orderService.ProcessOrder(OrderID)

	assert.Nil(t, err)
}

func TestOrderService_ProcessOrder_ErrorOrderAlreadyBeingProcessed(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockStore := orderMocks.NewMockOrderStorer(mockCtrl)
	mockStore.
		EXPECT().
		GetOrderByOrderID(OrderID).
		Return(&order.OrderDB{
			ID:      123,
			OrderID: OrderID,
			Status:  "PROCESSING",
		}, nil).
		Times(1)
	mockOrderProcessService := orderprocessMocks.NewMockOrderProcessor(mockCtrl)
	orderService := &order.OrderService{
		OrderStore:          mockStore,
		OrderProcessService: mockOrderProcessService,
	}

	err := orderService.ProcessOrder(OrderID)

	assert.Equal(t, order.OrderAlreadyBeingProcessedError, err)
}

func TestOrderService_ProcessOrder_ErrorOrderAlreadyFulfilled(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockStore := orderMocks.NewMockOrderStorer(mockCtrl)
	mockStore.
		EXPECT().
		GetOrderByOrderID(OrderID).
		Return(&order.OrderDB{
			ID:      123,
			OrderID: OrderID,
			Status:  "FULFILLED",
		}, nil).
		Times(1)
	mockOrderProcessService := orderprocessMocks.NewMockOrderProcessor(mockCtrl)
	orderService := &order.OrderService{
		OrderStore:          mockStore,
		OrderProcessService: mockOrderProcessService,
	}

	err := orderService.ProcessOrder(OrderID)

	assert.Equal(t, order.OrderAlreadyCompleteError, err)
}

func TestOrderService_ProcessOrder_ErrorOrderAlreadyFailed(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockStore := orderMocks.NewMockOrderStorer(mockCtrl)
	mockStore.
		EXPECT().
		GetOrderByOrderID(OrderID).
		Return(&order.OrderDB{
			ID:      123,
			OrderID: OrderID,
			Status:  "FAILED",
		}, nil).
		Times(1)
	mockOrderProcessService := orderprocessMocks.NewMockOrderProcessor(mockCtrl)
	orderService := &order.OrderService{
		OrderStore:          mockStore,
		OrderProcessService: mockOrderProcessService,
	}

	err := orderService.ProcessOrder(OrderID)

	assert.Equal(t, order.OrderAlreadyCompleteError, err)
}

func TestOrderService_ProcessOrder_ErrorOrderProcessServiceFailure(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockStore := orderMocks.NewMockOrderStorer(mockCtrl)
	mockStore.
		EXPECT().
		GetOrderByOrderID(OrderID).
		Return(&order.OrderDB{
			ID:      123,
			OrderID: OrderID,
			Status:  "CREATED",
		}, nil).
		Times(1)
	mockOrderProcessService := orderprocessMocks.NewMockOrderProcessor(mockCtrl)
	mockOrderProcessService.
		EXPECT().
		CreateOrder(OrderID).
		Return(errors.New("ERROR on order process service")).
		Times(1)
	orderService := &order.OrderService{
		OrderStore:          mockStore,
		OrderProcessService: mockOrderProcessService,
	}

	err := orderService.ProcessOrder(OrderID)

	assert.Equal(t, errs.GeneralError, err)
}

func TestOrderService_ProcessOrder_ErrorUpdateOrder(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockStore := orderMocks.NewMockOrderStorer(mockCtrl)
	mockStore.
		EXPECT().
		GetOrderByOrderID(OrderID).
		Return(&order.OrderDB{
			ID:      123,
			OrderID: OrderID,
			Status:  "CREATED",
		}, nil).
		Times(1)
	mockStore.
		EXPECT().
		UpdateOrder(&order.OrderDB{
			ID:      123,
			OrderID: OrderID,
			Status:  "PROCESSING",
		}).
		Return(0, errors.New("ERROR")).
		Times(1)
	mockOrderProcessService := orderprocessMocks.NewMockOrderProcessor(mockCtrl)
	mockOrderProcessService.
		EXPECT().
		CreateOrder(OrderID).
		Return(nil).
		Times(1)
	orderService := &order.OrderService{
		OrderStore:          mockStore,
		OrderProcessService: mockOrderProcessService,
	}

	err := orderService.ProcessOrder(OrderID)

	assert.Equal(t, errs.GeneralError, err)
}

func TestOrderService_ProcessOrder_ErrorOrderNotFound(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockStore := orderMocks.NewMockOrderStorer(mockCtrl)
	mockStore.
		EXPECT().
		GetOrderByOrderID(OrderID).
		Return(nil, sql.ErrNoRows).
		Times(1)
	mockOrderProcessService := orderprocessMocks.NewMockOrderProcessor(mockCtrl)
	orderService := &order.OrderService{
		OrderStore:          mockStore,
		OrderProcessService: mockOrderProcessService,
	}

	err := orderService.ProcessOrder(OrderID)

	assert.Equal(t, order.OrderNotFoundError, err)
}

func TestOrderService_CompleteOrder_SuccessWithStatusFulfilled(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockStore := orderMocks.NewMockOrderStorer(mockCtrl)
	mockStore.
		EXPECT().
		GetOrderByOrderID(OrderID).
		Return(&order.OrderDB{
			ID:      123,
			OrderID: OrderID,
			Status:  "PROCESSING",
		}, nil).
		Times(1)
	mockStore.
		EXPECT().
		UpdateOrder(&order.OrderDB{
			ID:      123,
			OrderID: OrderID,
			Status:  "FULFILLED",
		}).
		Return(123, nil).
		Times(1)
	mockOrderProcessService := orderprocessMocks.NewMockOrderProcessor(mockCtrl)
	orderService := &order.OrderService{
		OrderStore:          mockStore,
		OrderProcessService: mockOrderProcessService,
	}

	err := orderService.CompleteOrder(OrderID, order.ORDER_SUCCEEDED)

	assert.Nil(t, err)
}

func TestOrderService_CompleteOrder_SuccessWithStatusFailed(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockStore := orderMocks.NewMockOrderStorer(mockCtrl)
	mockStore.
		EXPECT().
		GetOrderByOrderID(OrderID).
		Return(&order.OrderDB{
			ID:      123,
			OrderID: OrderID,
			Status:  "PROCESSING",
		}, nil).
		Times(1)
	mockStore.
		EXPECT().
		UpdateOrder(&order.OrderDB{
			ID:      123,
			OrderID: OrderID,
			Status:  "FAILED",
		}).
		Return(123, nil).
		Times(1)
	mockOrderProcessService := orderprocessMocks.NewMockOrderProcessor(mockCtrl)
	orderService := &order.OrderService{
		OrderStore:          mockStore,
		OrderProcessService: mockOrderProcessService,
	}

	err := orderService.CompleteOrder(OrderID, order.ORDER_FAILED)

	assert.Nil(t, err)
}

func TestOrderService_CompleteOrder_ErrorOrderNotFound(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockStore := orderMocks.NewMockOrderStorer(mockCtrl)
	mockStore.
		EXPECT().
		GetOrderByOrderID(OrderID).
		Return(nil, sql.ErrNoRows).
		Times(1)
	mockOrderProcessService := orderprocessMocks.NewMockOrderProcessor(mockCtrl)
	orderService := &order.OrderService{
		OrderStore:          mockStore,
		OrderProcessService: mockOrderProcessService,
	}

	err := orderService.CompleteOrder(OrderID, order.ORDER_FAILED)

	assert.Equal(t, order.OrderNotFoundError, err)
}

func TestOrderService_CompleteOrder_ErrorInvalidOrderStatus(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockStore := orderMocks.NewMockOrderStorer(mockCtrl)
	mockStore.
		EXPECT().
		GetOrderByOrderID(OrderID).
		Return(&order.OrderDB{
			ID:      123,
			OrderID: OrderID,
			Status:  "PROCESSING",
		}, nil).
		Times(1)
	mockOrderProcessService := orderprocessMocks.NewMockOrderProcessor(mockCtrl)
	orderService := &order.OrderService{
		OrderStore:          mockStore,
		OrderProcessService: mockOrderProcessService,
	}

	err := orderService.CompleteOrder(OrderID, "invalidStatus")

	assert.Equal(t, order.UnknownProcessOrderStatusError, err)
}

func TestOrderService_CompleteOrder_ErrorGeneralError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockStore := orderMocks.NewMockOrderStorer(mockCtrl)
	mockStore.
		EXPECT().
		GetOrderByOrderID(OrderID).
		Return(nil, errors.New("ERROR")).
		Times(1)
	mockOrderProcessService := orderprocessMocks.NewMockOrderProcessor(mockCtrl)
	orderService := &order.OrderService{
		OrderStore:          mockStore,
		OrderProcessService: mockOrderProcessService,
	}

	err := orderService.CompleteOrder(OrderID, order.ORDER_FAILED)

	assert.Equal(t, errs.GeneralError, err)
}
