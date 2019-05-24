package order_test

import (
	"database/sql"
	"testing"

	"github.com/enrico5b1b4/order-service/order"
	"github.com/enrico5b1b4/order-service/test"
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestOrderStore_GetOrderByOrderID_Success(t *testing.T) {
	test.CheckSkipTest(t)
	orderID := uuid.NewV4()
	db := test.DBSetup(func(db *sqlx.DB) {
		test.InsertOrder(db, orderID, "CREATED")
	})
	orderStore := &order.OrderStore{DB: db}

	o, err := orderStore.GetOrderByOrderID(orderID)

	assert.Nil(t, err)
	assert.Equal(t, orderID, o.OrderID)
	assert.Equal(t, "CREATED", o.Status)
}

func TestOrderStore_GetOrderByOrderID_ErrorNoRows(t *testing.T) {
	test.CheckSkipTest(t)
	orderID := uuid.NewV4()
	db := test.DBSetup(func(db *sqlx.DB) {
		test.TruncateOrders(db)
	})
	orderStore := &order.OrderStore{DB: db}

	o, err := orderStore.GetOrderByOrderID(orderID)

	assert.Equal(t, sql.ErrNoRows, err)
	assert.Nil(t, o)
}

func TestOrderStore_GetOrderByID_Success(t *testing.T) {
	test.CheckSkipTest(t)
	orderID := uuid.NewV4()
	var id int
	db := test.DBSetup(func(db *sqlx.DB) {
		id = test.InsertOrder(db, orderID, "CREATED")
	})
	orderStore := &order.OrderStore{DB: db}

	o, err := orderStore.GetOrderByID(id)

	assert.Nil(t, err)
	assert.Equal(t, id, o.ID)
	assert.Equal(t, orderID, o.OrderID)
	assert.Equal(t, "CREATED", o.Status)
}

func TestOrderStore_GetOrderByID_ErrorNoRows(t *testing.T) {
	test.CheckSkipTest(t)
	db := test.DBSetup(func(db *sqlx.DB) {
		test.TruncateOrders(db)
	})
	orderStore := &order.OrderStore{DB: db}

	o, err := orderStore.GetOrderByID(1)

	assert.Equal(t, sql.ErrNoRows, err)
	assert.Nil(t, o)
}

func TestOrderStore_GetOrders_SuccessEmptyCollection(t *testing.T) {
	test.CheckSkipTest(t)
	db := test.DBSetup(func(db *sqlx.DB) {
		test.TruncateOrders(db)
	})
	orderStore := &order.OrderStore{DB: db}

	orders, err := orderStore.GetOrders("")

	assert.Nil(t, err)
	assert.Len(t, orders, 0)
}

func TestOrderStore_GetOrders_SuccessWithoutFilterStatus(t *testing.T) {
	test.CheckSkipTest(t)
	orderID1 := uuid.NewV4()
	orderID2 := uuid.NewV4()
	db := test.DBSetup(func(db *sqlx.DB) {
		test.TruncateOrders(db)
		test.InsertOrder(db, orderID1, "CREATED")
		test.InsertOrder(db, orderID2, "PROCESSING")
	})
	orderStore := &order.OrderStore{DB: db}

	orders, err := orderStore.GetOrders("")

	assert.Nil(t, err)
	assert.Len(t, orders, 2)
}

func TestOrderStore_GetOrders_SuccessWithFilterStatus(t *testing.T) {
	test.CheckSkipTest(t)
	orderID1 := uuid.NewV4()
	orderID2 := uuid.NewV4()
	db := test.DBSetup(func(db *sqlx.DB) {
		test.TruncateOrders(db)
		test.InsertOrder(db, orderID1, "CREATED")
		test.InsertOrder(db, orderID2, "PROCESSING")
	})
	orderStore := &order.OrderStore{DB: db}

	orders, err := orderStore.GetOrders("CREATED")

	assert.Nil(t, err)
	assert.Len(t, orders, 1)
	assert.Equal(t, orderID1, orders[0].OrderID)
	assert.Equal(t, "CREATED", orders[0].Status)
}

func TestOrderStore_CreateOrder_Success(t *testing.T) {
	test.CheckSkipTest(t)
	orderID := uuid.NewV4()
	oDB := &order.OrderDB{
		OrderID: orderID,
		Status:  "CREATED",
	}
	db := test.DBSetup(func(db *sqlx.DB) {
		test.TruncateOrders(db)
	})
	orderStore := &order.OrderStore{DB: db}

	oDBID, err := orderStore.CreateOrder(oDB)

	assert.Nil(t, err)
	assert.NotEqual(t, 0, oDBID)

	o, err := orderStore.GetOrderByOrderID(orderID)

	assert.Nil(t, err)
	assert.Equal(t, orderID, o.OrderID)
	assert.Equal(t, "CREATED", o.Status)
}

func TestOrderStore_CreateOrder_ErrorUniqueViolationError(t *testing.T) {
	test.CheckSkipTest(t)
	orderID := uuid.NewV4()
	oDB := &order.OrderDB{
		OrderID: orderID,
		Status:  "CREATED",
	}
	db := test.DBSetup(func(db *sqlx.DB) {
		test.TruncateOrders(db)
		test.InsertOrder(db, orderID, "CREATED")
	})
	orderStore := &order.OrderStore{DB: db}

	oDBID, err := orderStore.CreateOrder(oDB)

	assert.Equal(t, err, order.DBOrderAlreadyExistsError)
	assert.Equal(t, 0, oDBID)
}

func TestOrderStore_UpdateOrder_Success(t *testing.T) {
	test.CheckSkipTest(t)
	orderID := uuid.NewV4()
	var id int
	db := test.DBSetup(func(db *sqlx.DB) {
		test.TruncateOrders(db)
		id = test.InsertOrder(db, orderID, "CREATED")
	})
	oDB := &order.OrderDB{
		ID:      id,
		OrderID: orderID,
		Status:  "PROCESSING",
	}
	orderStore := &order.OrderStore{DB: db}

	oDBID, err := orderStore.UpdateOrder(oDB)

	assert.Nil(t, err)
	assert.NotEqual(t, 0, oDBID)

	o, err := orderStore.GetOrderByOrderID(orderID)

	assert.Nil(t, err)
	assert.Equal(t, orderID, o.OrderID)
	assert.Equal(t, "PROCESSING", o.Status)
}

func TestOrderStore_UpdateOrder_ErrorNoRows(t *testing.T) {
	test.CheckSkipTest(t)
	orderID := uuid.NewV4()
	db := test.DBSetup(func(db *sqlx.DB) {
		test.TruncateOrders(db)
	})
	oDB := &order.OrderDB{
		ID:      1,
		OrderID: orderID,
		Status:  "PROCESSING",
	}
	orderStore := &order.OrderStore{DB: db}

	oDBID, err := orderStore.UpdateOrder(oDB)

	assert.Equal(t, sql.ErrNoRows, err)
	assert.Equal(t, 0, oDBID)
}
