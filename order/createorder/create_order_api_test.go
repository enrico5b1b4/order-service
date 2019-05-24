package createorder_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/enrico5b1b4/order-service/test"
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
)

func TestAPICreateOrder_Success(t *testing.T) {
	test.CheckSkipTest(t)
	ID := uuid.NewV4()

	test.ApiTest().
		Post("/order").
		JSON(fmt.Sprintf(`{"id": "%s"}`, ID.String())).
		Expect(t).
		Status(http.StatusOK).
		Body(fmt.Sprintf(`{"order": {"id": "%s", "status": "CREATED"}}`, ID.String())).
		End()
}

func TestAPICreateOrder_InvalidRequestBody(t *testing.T) {
	test.CheckSkipTest(t)
	test.ApiTest().
		Post("/order").
		JSON(`{"NOT_A_VALID_KEY": ""}`).
		Expect(t).
		Status(http.StatusBadRequest).
		End()
}

func TestAPICreateOrder_OrderAlreadyExists(t *testing.T) {
	test.CheckSkipTest(t)
	ID := uuid.NewV4()

	test.DBSetup(func(db *sqlx.DB) {
		test.InsertOrder(db, ID, "CREATED")
	})

	test.ApiTest().
		Post("/order").
		JSON(fmt.Sprintf(`{"id": "%s"}`, ID.String())).
		Expect(t).
		Status(http.StatusBadRequest).
		Body(`{"code":"5","message":"Order id already exists."}`).
		End()
}
