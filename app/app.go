package app

import (
	"net/http"
	"time"

	"github.com/enrico5b1b4/order-service/order/processorder"

	"github.com/enrico5b1b4/order-service/order"
	"github.com/enrico5b1b4/order-service/order/completeorder"
	"github.com/enrico5b1b4/order-service/order/createorder"
	"github.com/enrico5b1b4/order-service/order/getorder"
	"github.com/enrico5b1b4/order-service/order/getorders"
	"github.com/enrico5b1b4/order-service/orderprocess"
	"github.com/enrico5b1b4/order-service/request"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func New(DB *sqlx.DB, orderProcessServiceURL, completeOrderCallbackURL string) *echo.Echo {
	// init services/stores/clients
	orderProcessClient := &orderprocess.OrderProcessClient{
		Client: &http.Client{
			Timeout: time.Second * 5,
		},
		BaseURL:                   orderProcessServiceURL,
		OrderProcessedCallbackURL: completeOrderCallbackURL,
	}
	orderStore := &order.OrderStore{
		DB: DB,
	}
	orderService := &order.OrderService{
		OrderStore:          orderStore,
		OrderProcessService: orderProcessClient,
	}

	// init app/routes
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.HTTPErrorHandler = CustomHTTPErrorHandler
	e.Validator = &request.Validator{Validator: request.NewRequestValidator()}

	e.GET("/ping", ping)
	e.GET("/order", getorders.GetOrders(orderService))
	e.GET("/order/:orderID", getorder.GetOrder(orderService))
	e.POST("/order", createorder.CreateOrder(orderService))
	e.POST("/process_order", processorder.ProcessOrder(orderService))
	e.POST("/complete_order", completeorder.CompleteOrder(orderService))

	return e
}

func ping(c echo.Context) error {
	return c.String(http.StatusOK, "")
}
