package getorder_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/enrico5b1b4/order-service/app"
	"github.com/enrico5b1b4/order-service/errors"
	"github.com/enrico5b1b4/order-service/order"
	"github.com/enrico5b1b4/order-service/order/getorder"
	"github.com/enrico5b1b4/order-service/order/mocks"
	"github.com/enrico5b1b4/order-service/test"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

var ID = uuid.Must(uuid.FromString("f6c7d890-ca9b-4147-aa5b-2a41cb17f95a"))

func TestGetOrderHandler_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockService := mocks.NewMockOrderServicer(mockCtrl)
	mockService.EXPECT().GetOrderByID(ID).Return(&order.Order{
		ID:     ID,
		Status: order.FULFILLED,
	}, nil).Times(1)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := test.NewEchoTestContext(req, rec)
	c.SetPath("/order/:orderID")
	c.SetParamNames("orderID")
	c.SetParamValues(ID.String())

	if assert.NoError(t, getorder.GetOrder(mockService)(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `{"order": {"id": "f6c7d890-ca9b-4147-aa5b-2a41cb17f95a", "status": "FULFILLED"}}`, rec.Body.String())
	}
}

func TestGetOrderHandler_NotFoundError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockService := mocks.NewMockOrderServicer(mockCtrl)
	mockService.EXPECT().GetOrderByID(ID).Return(nil, order.OrderNotFoundError).Times(1)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := test.NewEchoTestContext(req, rec)
	c.SetPath("/order/:orderID")
	c.SetParamNames("orderID")
	c.SetParamValues(ID.String())

	app.CustomHTTPErrorHandler(getorder.GetOrder(mockService)(c), c)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.JSONEq(t, `{"code": "3", "message":"Order not found."}`, rec.Body.String())
}

func TestGetOrderHandler_GeneralError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockService := mocks.NewMockOrderServicer(mockCtrl)
	mockService.EXPECT().GetOrderByID(ID).Return(nil, errors.GeneralError).Times(1)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := test.NewEchoTestContext(req, rec)
	c.SetPath("/order/:orderID")
	c.SetParamNames("orderID")
	c.SetParamValues(ID.String())

	app.CustomHTTPErrorHandler(getorder.GetOrder(mockService)(c), c)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.JSONEq(t, `{"code": "0", "message":"Something went wrong."}`, rec.Body.String())
}
