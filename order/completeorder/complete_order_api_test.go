package completeorder_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/enrico5b1b4/order-service/test"
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
)

func TestAPICompleteOrder_Success(t *testing.T) {
	test.CheckSkipTest(t)
	ID := uuid.NewV4()

	test.DBSetup(func(db *sqlx.DB) {
		test.InsertOrder(db, ID, "PROCESSING")
	})

	test.ApiTest().
		Post("/complete_order").
		JSON(fmt.Sprintf(`{"order_id": "%s", "status": "FAILED"}`, ID.String())).
		Expect(t).
		Status(http.StatusOK).
		End()

	test.ApiTest().
		Get(fmt.Sprintf("/order/%s", ID.String())).
		Expect(t).
		Status(http.StatusOK).
		Body(fmt.Sprintf(`{"order": {"id": "%s", "status": "FAILED"}}`, ID.String())).
		End()
}

func TestAPICompleteOrder_FailureOrderNotFound(t *testing.T) {
	test.CheckSkipTest(t)
	ID := uuid.NewV4()

	test.DBSetup(func(db *sqlx.DB) {
		test.TruncateOrders(db)
	})

	test.ApiTest().
		Post("/complete_order").
		JSON(fmt.Sprintf(`{"order_id": "%s", "status": "FAILED"}`, ID.String())).
		Expect(t).
		Status(http.StatusNotFound).
		Body(`{"code":"3","message":"Order not found."}`).
		End()
}

func TestAPICompleteOrder_FailureInvalidStatus(t *testing.T) {
	test.CheckSkipTest(t)
	ID := uuid.NewV4()

	test.DBSetup(func(db *sqlx.DB) {
		test.TruncateOrders(db)
		test.InsertOrder(db, ID, "PROCESSING")
	})

	test.ApiTest().
		Post("/complete_order").
		JSON(fmt.Sprintf(`{"order_id": "%s", "status": "NOT_A_VALID_STATUS"}`, ID.String())).
		Expect(t).
		Status(http.StatusBadRequest).
		End()
}
