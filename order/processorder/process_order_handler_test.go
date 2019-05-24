package processorder_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/enrico5b1b4/order-service/app"
	"github.com/enrico5b1b4/order-service/errors"
	"github.com/enrico5b1b4/order-service/order"
	"github.com/enrico5b1b4/order-service/order/mocks"
	"github.com/enrico5b1b4/order-service/order/processorder"
	"github.com/enrico5b1b4/order-service/test"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

var ID = uuid.Must(uuid.FromString("f6c7d890-ca9b-4147-aa5b-2a41cb17f95a"))

func TestProcessOrderHandler_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockService := mocks.NewMockOrderServicer(mockCtrl)
	mockService.EXPECT().ProcessOrder(ID).Return(nil).Times(1)

	req := httptest.NewRequest(
		http.MethodPost,
		"/",
		strings.NewReader(`{"id": "f6c7d890-ca9b-4147-aa5b-2a41cb17f95a"}`),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := test.NewEchoTestContext(req, rec)

	if assert.NoError(t, processorder.ProcessOrder(mockService)(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestProcessOrderHandler_FailureOrderAlreadyBeingProcessed(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockService := mocks.NewMockOrderServicer(mockCtrl)
	mockService.EXPECT().ProcessOrder(ID).Return(order.OrderAlreadyBeingProcessedError).Times(1)

	req := httptest.NewRequest(
		http.MethodPost,
		"/",
		strings.NewReader(`{"id": "f6c7d890-ca9b-4147-aa5b-2a41cb17f95a"}`),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := test.NewEchoTestContext(req, rec)

	app.CustomHTTPErrorHandler(processorder.ProcessOrder(mockService)(c), c)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, `{"code": "7", "message":"Cannot start processing order since it is already being processed."}`, rec.Body.String())
}

func TestProcessOrderHandler_FailureOrderAlreadyCompleteError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockService := mocks.NewMockOrderServicer(mockCtrl)
	mockService.EXPECT().ProcessOrder(ID).Return(order.OrderAlreadyCompleteError).Times(1)

	req := httptest.NewRequest(
		http.MethodPost,
		"/",
		strings.NewReader(`{"id": "f6c7d890-ca9b-4147-aa5b-2a41cb17f95a"}`),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := test.NewEchoTestContext(req, rec)

	app.CustomHTTPErrorHandler(processorder.ProcessOrder(mockService)(c), c)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, `{"code":"8", "message":"Cannot process a completed order."}`, rec.Body.String())
}

func TestProcessOrderHandler_FailureGeneralError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockService := mocks.NewMockOrderServicer(mockCtrl)
	mockService.EXPECT().ProcessOrder(ID).Return(errors.GeneralError).Times(1)

	req := httptest.NewRequest(
		http.MethodPost,
		"/",
		strings.NewReader(`{"id": "f6c7d890-ca9b-4147-aa5b-2a41cb17f95a"}`),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := test.NewEchoTestContext(req, rec)

	app.CustomHTTPErrorHandler(processorder.ProcessOrder(mockService)(c), c)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.JSONEq(t, `{"code":"0", "message":"Something went wrong."}`, rec.Body.String())
}

func TestProcessOrderHandler_FailureInvalidRequest(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockService := mocks.NewMockOrderServicer(mockCtrl)

	req := httptest.NewRequest(
		http.MethodPost,
		"/",
		strings.NewReader(`{"NOT_A_VALID_KEY": "f6c7d890-ca9b-4147-aa5b-2a41cb17f95a"}`),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := test.NewEchoTestContext(req, rec)

	app.CustomHTTPErrorHandler(processorder.ProcessOrder(mockService)(c), c)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
