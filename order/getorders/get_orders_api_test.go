package getorders_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/enrico5b1b4/order-service/test"
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
)

func TestAPIGetOrders_SuccessWithOrders(t *testing.T) {
	test.CheckSkipTest(t)
	ID1 := uuid.NewV4()
	ID2 := uuid.NewV4()

	test.DBSetup(func(db *sqlx.DB) {
		test.TruncateOrders(db)
		test.InsertOrder(db, ID1, "CREATED")
		test.InsertOrder(db, ID2, "FULFILLED")
	})

	test.ApiTest().
		Get("/order").
		Expect(t).
		Status(http.StatusOK).
		Body(fmt.Sprintf(
			`{"orders":[{"id":"%s","status":"CREATED"},{"id":"%s","status":"FULFILLED"}]}`,
			ID1.String(),
			ID2.String(),
		)).
		End()
}

func TestAPIGetOrders_SuccessWithFilteredOrders(t *testing.T) {
	test.CheckSkipTest(t)
	ID1 := uuid.NewV4()
	ID2 := uuid.NewV4()

	test.DBSetup(func(db *sqlx.DB) {
		test.TruncateOrders(db)
		test.InsertOrder(db, ID1, "CREATED")
		test.InsertOrder(db, ID2, "FULFILLED")
	})

	test.ApiTest().
		Get("/order").
		Query("status", "CREATED").
		Expect(t).
		Status(http.StatusOK).
		Body(fmt.Sprintf(`{"orders":[{"id":"%s","status":"CREATED"}]}`, ID1.String())).
		End()
}

func TestAPIGetOrders_SuccessWithNoOrders(t *testing.T) {
	test.CheckSkipTest(t)
	test.DBSetup(func(db *sqlx.DB) {
		test.TruncateOrders(db)
	})

	test.ApiTest().
		Get("/order").
		Expect(t).
		Status(http.StatusOK).
		Body(`{"orders":[]}`).
		End()
}
