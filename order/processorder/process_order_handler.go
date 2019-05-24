package processorder

import (
	"net/http"

	"github.com/enrico5b1b4/order-service/order"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

type ProcessOrderRequest struct {
	ID uuid.UUID `json:"id" validate:"required"`
}

func ProcessOrder(service order.OrderServicer) func(c echo.Context) error {
	return func(c echo.Context) error {
		poRequest := new(ProcessOrderRequest)
		if err := c.Bind(poRequest); err != nil {
			return err
		}
		if err := c.Validate(poRequest); err != nil {
			return err
		}

		createOrderErr := service.ProcessOrder(poRequest.ID)
		if createOrderErr != nil {
			return createOrderErr
		}

		return c.JSON(http.StatusOK, ``)
	}
}
