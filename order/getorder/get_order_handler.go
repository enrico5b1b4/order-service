package getorder

import (
	"net/http"

	"github.com/enrico5b1b4/order-service/order"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

type OrderResponse struct {
	Order *order.Order `json:"order"`
}

func GetOrder(service order.OrderServicer) func(c echo.Context) error {
	return func(c echo.Context) error {
		orderID, err := uuid.FromString(c.Param("orderID"))
		if err != nil {
			return err
		}

		o, getOrderErr := service.GetOrderByID(orderID)
		if getOrderErr != nil {
			return getOrderErr
		}

		return c.JSON(http.StatusOK, OrderResponse{Order: o})
	}
}
