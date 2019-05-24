package createorder

import (
	"net/http"

	"github.com/enrico5b1b4/order-service/order"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

type CreateOrderRequest struct {
	ID uuid.UUID `json:"id" validate:"required"`
}

type CreateOrderResponse struct {
	Order *order.Order `json:"order"`
}

func CreateOrder(service order.OrderServicer) func(c echo.Context) error {
	return func(c echo.Context) error {
		createOrderRequest := new(CreateOrderRequest)
		if err := c.Bind(createOrderRequest); err != nil {
			return err
		}
		if err := c.Validate(createOrderRequest); err != nil {
			return err
		}

		o, createOrderErr := service.CreateOrder(mapToOrder(createOrderRequest))
		if createOrderErr != nil {
			return createOrderErr
		}

		return c.JSON(http.StatusOK, CreateOrderResponse{Order: o})
	}
}

func mapToOrder(createOrderRequest *CreateOrderRequest) *order.Order {
	return &order.Order{
		ID: createOrderRequest.ID,
	}
}
