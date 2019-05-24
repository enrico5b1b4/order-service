package processorder_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/enrico5b1b4/order-service/test"
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
	"github.com/steinfletcher/apitest"
)

func TestAPIProcessOrder_Success(t *testing.T) {
	test.CheckSkipTest(t)
	ID := uuid.NewV4()

	test.DBSetup(func(db *sqlx.DB) {
		test.InsertOrder(db, ID, "CREATED")
	})

	test.ApiTest().
		Mocks(mockPostOrderProcessSuccess(ID)).
		Post("/process_order").
		JSON(fmt.Sprintf(`{"id": "%s"}`, ID.String())).
		Expect(t).
		Status(http.StatusOK).
		End()
}

func TestAPIProcessOrder_FailureOrderAlreadyComplete(t *testing.T) {
	test.CheckSkipTest(t)
	ID := uuid.NewV4()

	test.DBSetup(func(db *sqlx.DB) {
		test.InsertOrder(db, ID, "FAILED")
	})

	test.ApiTest().
		Post("/process_order").
		JSON(fmt.Sprintf(`{"id": "%s"}`, ID.String())).
		Expect(t).
		Status(http.StatusBadRequest).
		Body(`{"code":"8", "message":"Cannot process a completed order."}`).
		End()
}

func TestAPIProcessOrder_FailureOrderAlreadyBeingProcessed(t *testing.T) {
	test.CheckSkipTest(t)
	ID := uuid.NewV4()

	test.DBSetup(func(db *sqlx.DB) {
		test.InsertOrder(db, ID, "PROCESSING")
	})

	test.ApiTest().
		Post("/process_order").
		JSON(fmt.Sprintf(`{"id": "%s"}`, ID.String())).
		Expect(t).
		Status(http.StatusBadRequest).
		Body(`{"code":"7", "message":"Cannot start processing order since it is already being processed."}`).
		End()
}

func TestAPIProcessOrder_FailureOrderProcessAlreadyExists(t *testing.T) {
	test.CheckSkipTest(t)
	ID := uuid.NewV4()

	test.DBSetup(func(db *sqlx.DB) {
		test.InsertOrder(db, ID, "CREATED")
	})

	test.ApiTest().
		Mocks(mockPostOrderProcessFailureAlreadyExists(ID)).
		Post("/process_order").
		JSON(fmt.Sprintf(`{"id": "%s"}`, ID.String())).
		Expect(t).
		Status(http.StatusInternalServerError).
		Body(`{"code":"0", "message":"Something went wrong."}`).
		End()
}

func TestAPIProcessOrder_FailureOrderProcessGeneralError(t *testing.T) {
	test.CheckSkipTest(t)
	ID := uuid.NewV4()

	test.DBSetup(func(db *sqlx.DB) {
		test.InsertOrder(db, ID, "CREATED")
	})

	test.ApiTest().
		Mocks(mockPostOrderProcessFailureGeneralError(ID)).
		Post("/process_order").
		JSON(fmt.Sprintf(`{"id": "%s"}`, ID.String())).
		Expect(t).
		Status(http.StatusInternalServerError).
		Body(`{"code":"0", "message":"Something went wrong."}`).
		End()
}

func mockPostOrderProcessSuccess(ID uuid.UUID) *apitest.Mock {
	return apitest.NewMock().
		Post("http://orderprocessservice.local/order_process").
		Query("callback_url", "http://localhost/complete_order").
		Header("Content-Type", "application/json").
		Body(fmt.Sprintf(`{"order_id":"%s"}`, ID.String())).
		RespondWith().
		Status(http.StatusAccepted).
		End()
}

func mockPostOrderProcessFailureAlreadyExists(ID uuid.UUID) *apitest.Mock {
	return apitest.NewMock().
		Post("http://orderprocessservice.local/order_process").
		Query("callback_url", "http://localhost/complete_order").
		Header("Content-Type", "application/json").
		Body(fmt.Sprintf(`{"order_id":"%s"}`, ID.String())).
		RespondWith().
		Status(http.StatusBadRequest).
		End()
}

func mockPostOrderProcessFailureGeneralError(ID uuid.UUID) *apitest.Mock {
	return apitest.NewMock().
		Post("http://orderprocessservice.local/order_process").
		Query("callback_url", "http://localhost/complete_order").
		Header("Content-Type", "application/json").
		Body(fmt.Sprintf(`{"order_id":"%s"}`, ID.String())).
		RespondWith().
		Status(http.StatusInternalServerError).
		End()
}
