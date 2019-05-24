package completeorder_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/enrico5b1b4/order-service/app"
	"github.com/enrico5b1b4/order-service/errors"
	"github.com/enrico5b1b4/order-service/order"
	"github.com/enrico5b1b4/order-service/order/completeorder"
	"github.com/enrico5b1b4/order-service/order/mocks"
	"github.com/enrico5b1b4/order-service/test"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

var ID = uuid.Must(uuid.FromString("f6c7d890-ca9b-4147-aa5b-2a41cb17f95a"))

func TestCompleteOrderHandler_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockService := mocks.NewMockOrderServicer(mockCtrl)
	mockService.EXPECT().CompleteOrder(ID, order.ORDER_SUCCEEDED).Return(nil).Times(1)

	req := newCompleteOrderTestRequest(ID.String(), "SUCCEEDED")
	rec := httptest.NewRecorder()
	c := test.NewEchoTestContext(req, rec)

	if assert.NoError(t, completeorder.CompleteOrder(mockService)(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestCompleteOrderHandler_FailureInvalidRequest(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockService := mocks.NewMockOrderServicer(mockCtrl)

	req := newCompleteOrderTestRequest(ID.String(), "NOT_A_VALID_STATUS")
	rec := httptest.NewRecorder()
	c := test.NewEchoTestContext(req, rec)

	app.CustomHTTPErrorHandler(completeorder.CompleteOrder(mockService)(c), c)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCompleteOrderHandler_FailureOrderNotFound(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockService := mocks.NewMockOrderServicer(mockCtrl)
	mockService.EXPECT().CompleteOrder(ID, order.ORDER_SUCCEEDED).Return(order.OrderNotFoundError).Times(1)

	req := newCompleteOrderTestRequest(ID.String(), "SUCCEEDED")
	rec := httptest.NewRecorder()
	c := test.NewEchoTestContext(req, rec)

	app.CustomHTTPErrorHandler(completeorder.CompleteOrder(mockService)(c), c)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.JSONEq(t, `{"code": "3", "message":"Order not found."}`, rec.Body.String())
}

func TestCompleteOrderHandler_GeneralError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockService := mocks.NewMockOrderServicer(mockCtrl)
	mockService.EXPECT().CompleteOrder(ID, order.ORDER_SUCCEEDED).Return(errors.GeneralError).Times(1)

	req := newCompleteOrderTestRequest(ID.String(), "SUCCEEDED")
	rec := httptest.NewRecorder()
	c := test.NewEchoTestContext(req, rec)

	app.CustomHTTPErrorHandler(completeorder.CompleteOrder(mockService)(c), c)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.JSONEq(t, `{"code": "0", "message":"Something went wrong."}`, rec.Body.String())
}

func newCompleteOrderTestRequest(orderID, status string) *http.Request {
	req := httptest.NewRequest(
		http.MethodPost,
		"/",
		strings.NewReader(fmt.Sprintf(`{"order_id": "%s", "status": "%s"}`, orderID, status)),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	return req
}
