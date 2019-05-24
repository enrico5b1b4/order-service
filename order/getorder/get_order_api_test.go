package getorder_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/enrico5b1b4/order-service/test"
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
)

func TestAPIGetOrder_Success(t *testing.T) {
	test.CheckSkipTest(t)
	ID := uuid.NewV4()

	test.DBSetup(func(db *sqlx.DB) {
		test.InsertOrder(db, ID, "CREATED")
	})

	test.ApiTest().
		Get(fmt.Sprintf("/order/%s", ID.String())).
		Expect(t).
		Status(http.StatusOK).
		Body(fmt.Sprintf(`{"order": {"id": "%s", "status": "CREATED"}}`, ID.String())).
		End()
}

func TestAPIGetOrder_OrderNotFound(t *testing.T) {
	test.CheckSkipTest(t)
	ID := uuid.NewV4()

	test.ApiTest().
		Get(fmt.Sprintf("/order/%s", ID.String())).
		Expect(t).
		Status(http.StatusNotFound).
		Body(`{"code":"3", "message":"Order not found."}`).
		End()
}
