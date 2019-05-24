package createorder_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/enrico5b1b4/order-service/app"
	"github.com/enrico5b1b4/order-service/errors"
	"github.com/enrico5b1b4/order-service/order"
	"github.com/enrico5b1b4/order-service/order/createorder"
	"github.com/enrico5b1b4/order-service/order/mocks"
	"github.com/enrico5b1b4/order-service/test"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

var ID = uuid.Must(uuid.FromString("f6c7d890-ca9b-4147-aa5b-2a41cb17f95a"))

func TestCreateOrderHandler_Success(t *testing.T) {
	newOrder := &order.Order{ID: ID}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockService := mocks.NewMockOrderServicer(mockCtrl)
	mockService.EXPECT().CreateOrder(newOrder).Return(&order.Order{
		ID:     ID,
		Status: order.CREATED,
	}, nil).Times(1)

	req := newCreateOrderTestRequest(ID.String())
	rec := httptest.NewRecorder()
	c := test.NewEchoTestContext(req, rec)

	if assert.NoError(t, createorder.CreateOrder(mockService)(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `{"order": {"id": "f6c7d890-ca9b-4147-aa5b-2a41cb17f95a", "status": "CREATED"}}`, rec.Body.String())
	}
}

func TestCreateOrderHandler_GeneralError(t *testing.T) {
	newOrder := &order.Order{ID: ID}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockService := mocks.NewMockOrderServicer(mockCtrl)
	mockService.EXPECT().CreateOrder(newOrder).Return(nil, errors.GeneralError).Times(1)

	req := newCreateOrderTestRequest(ID.String())
	rec := httptest.NewRecorder()
	c := test.NewEchoTestContext(req, rec)

	app.CustomHTTPErrorHandler(createorder.CreateOrder(mockService)(c), c)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.JSONEq(t, `{"code": "0", "message":"Something went wrong."}`, rec.Body.String())
}

func TestCreateOrderHandler_OrderAlreadyExistsError(t *testing.T) {
	newOrder := &order.Order{ID: ID}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockService := mocks.NewMockOrderServicer(mockCtrl)
	mockService.EXPECT().CreateOrder(newOrder).Return(nil, order.OrderAlreadyExistsError).Times(1)

	req := newCreateOrderTestRequest(ID.String())
	rec := httptest.NewRecorder()
	c := test.NewEchoTestContext(req, rec)

	app.CustomHTTPErrorHandler(createorder.CreateOrder(mockService)(c), c)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, `{"code": "5", "message":"Order id already exists."}`, rec.Body.String())
}

func newCreateOrderTestRequest(orderID string) *http.Request {
	req := httptest.NewRequest(
		http.MethodPost,
		"/",
		strings.NewReader(fmt.Sprintf(`{"id": "%s"}`, orderID)),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	return req
}
