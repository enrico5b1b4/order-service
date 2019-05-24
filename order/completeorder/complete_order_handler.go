package completeorder

import (
	"net/http"

	"github.com/enrico5b1b4/order-service/order"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/go-playground/validator.v9"
)

type CompleteOrderRequest struct {
	ID                 uuid.UUID `json:"order_id" validate:"required"`
	ProcessOrderStatus string    `json:"status" validate:"required,completeOrderStatusValidation"`
}

func CompleteOrder(service order.OrderServicer) func(c echo.Context) error {
	return func(c echo.Context) error {
		coRequest := new(CompleteOrderRequest)
		if err := c.Bind(coRequest); err != nil {
			return err
		}
		if err := c.Validate(coRequest); err != nil {
			return err
		}

		createOrderErr := service.CompleteOrder(coRequest.ID, coRequest.ProcessOrderStatus)
		if createOrderErr != nil {
			return createOrderErr
		}

		return c.JSON(http.StatusOK, ``)
	}
}

func CompleteOrderStatusValidation(fl validator.FieldLevel) bool {
	return fl.Field().String() == order.ORDER_SUCCEEDED || fl.Field().String() == order.ORDER_FAILED
}
