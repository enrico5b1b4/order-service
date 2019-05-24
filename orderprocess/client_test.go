package orderprocess_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/enrico5b1b4/order-service/errors"
	"github.com/enrico5b1b4/order-service/orderprocess"
	uuid "github.com/satori/go.uuid"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/assert"
)

func TestOrderProcessClient_CreateOrder_Success(t *testing.T) {
	ID := uuid.NewV4()
	client := orderprocess.OrderProcessClient{
		BaseURL:                   "http://orderprocessservice.local",
		Client:                    &http.Client{},
		OrderProcessedCallbackURL: "http://localhost/complete_order",
	}
	defer apitest.NewMock().
		Post("http://orderprocessservice.local/order_process").
		Query("callback_url", "http://localhost/complete_order").
		Header("Content-Type", "application/json").
		Body(fmt.Sprintf(`{"order_id":"%s"}`, ID.String())).
		RespondWith().
		Status(http.StatusAccepted).
		EndStandalone()()

	err := client.CreateOrder(ID)

	assert.Nil(t, err)
}

func TestOrderProcessClient_CreateOrder_Error(t *testing.T) {
	ID := uuid.NewV4()
	client := orderprocess.OrderProcessClient{
		BaseURL:                   "http://orderprocessservice.local",
		Client:                    &http.Client{},
		OrderProcessedCallbackURL: "http://localhost/complete_order",
	}
	defer apitest.NewMock().
		Post("http://orderprocessservice.local/order_process").
		Query("callback_url", "http://localhost/complete_order").
		Header("Content-Type", "application/json").
		Body(fmt.Sprintf(`{"order_id":"%s"}`, ID.String())).
		RespondWith().
		Status(http.StatusBadRequest).
		EndStandalone()()

	err := client.CreateOrder(ID)

	assert.Error(t, errors.New("orderprocess: error creating order"), err)
}
