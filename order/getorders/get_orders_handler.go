package getorders

import (
	"net/http"

	"github.com/enrico5b1b4/order-service/order"
	"github.com/labstack/echo/v4"
)

type OrdersResponse struct {
	Orders []*order.Order `json:"orders"`
}

func GetOrders(service order.OrderServicer) func(c echo.Context) error {
	return func(c echo.Context) error {
		status := c.QueryParam("status")

		orders, getOrdersErr := service.GetOrders(status)
		if getOrdersErr != nil {
			return getOrdersErr
		}

		return c.JSON(http.StatusOK, OrdersResponse{Orders: orders})
	}
}
