package getorders_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/enrico5b1b4/order-service/app"
	"github.com/enrico5b1b4/order-service/errors"
	"github.com/enrico5b1b4/order-service/order"
	"github.com/enrico5b1b4/order-service/order/getorders"
	"github.com/enrico5b1b4/order-service/order/mocks"
	"github.com/enrico5b1b4/order-service/test"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

var ID = uuid.Must(uuid.FromString("f6c7d890-ca9b-4147-aa5b-2a41cb17f95a"))

func TestGetOrdersHandler_SuccessWithoutFilteredStatus(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockService := mocks.NewMockOrderServicer(mockCtrl)
	mockService.EXPECT().GetOrders("").Return([]*order.Order{{
		ID:     ID,
		Status: order.FULFILLED,
	}}, nil).Times(1)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := test.NewEchoTestContext(req, rec)

	if assert.NoError(t, getorders.GetOrders(mockService)(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `{"orders": [{"id": "f6c7d890-ca9b-4147-aa5b-2a41cb17f95a", "status": "FULFILLED"}]}`, rec.Body.String())
	}
}

func TestGetOrdersHandler_SuccessWithFilteredStatus(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockService := mocks.NewMockOrderServicer(mockCtrl)
	mockService.EXPECT().GetOrders("FULFILLED").Return([]*order.Order{{
		ID:     ID,
		Status: order.FULFILLED,
	}}, nil).Times(1)

	q := make(url.Values)
	q.Set("status", "FULFILLED")
	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := test.NewEchoTestContext(req, rec)

	if assert.NoError(t, getorders.GetOrders(mockService)(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `{"orders": [{"id": "f6c7d890-ca9b-4147-aa5b-2a41cb17f95a", "status": "FULFILLED"}]}`, rec.Body.String())
	}
}

func TestGetOrdersHandler_SuccessNoOrders(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockService := mocks.NewMockOrderServicer(mockCtrl)
	mockService.EXPECT().GetOrders("").Return([]*order.Order{}, nil).Times(1)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := test.NewEchoTestContext(req, rec)

	if assert.NoError(t, getorders.GetOrders(mockService)(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `{"orders": []}`, rec.Body.String())
	}
}

func TestGetOrdersHandler_GeneralError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockService := mocks.NewMockOrderServicer(mockCtrl)
	mockService.EXPECT().GetOrders("").Return(nil, errors.GeneralError).Times(1)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := test.NewEchoTestContext(req, rec)

	app.CustomHTTPErrorHandler(getorders.GetOrders(mockService)(c), c)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.JSONEq(t, `{"code": "0", "message":"Something went wrong."}`, rec.Body.String())
}
