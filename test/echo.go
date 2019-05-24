package test

import (
	"net/http"

	"github.com/enrico5b1b4/order-service/app"
	"github.com/enrico5b1b4/order-service/request"
	"github.com/labstack/echo/v4"
)

func NewEchoTestContext(r *http.Request, w http.ResponseWriter) echo.Context {
	e := echo.New()
	e.Validator = &request.Validator{Validator: request.NewRequestValidator()}
	e.HTTPErrorHandler = app.CustomHTTPErrorHandler
	return e.NewContext(r, w)
}
