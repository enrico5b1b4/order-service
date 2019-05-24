package orderprocess

//go:generate mockgen -destination=./mocks/mock_OrderProcessor.go -package=mocks github.com/enrico5b1b4/order-service/orderprocess OrderProcessor

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	uuid "github.com/satori/go.uuid"
)

type OrderProcessor interface {
	CreateOrder(uuid.UUID) error
}

type OrderProcessClient struct {
	BaseURL                   string
	Client                    *http.Client
	OrderProcessedCallbackURL string
}

func (c *OrderProcessClient) CreateOrder(ID uuid.UUID) error {
	apiUrl, err := url.Parse(fmt.Sprintf("%s/order_process", c.BaseURL))
	if err != nil {
		return err
	}

	q := apiUrl.Query()
	q.Add("callback_url", c.OrderProcessedCallbackURL)
	apiUrl.RawQuery = q.Encode()

	o := &Order{ID: ID}
	requestBody, err := json.Marshal(o)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, apiUrl.String(), bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return errors.New("orderprocess: error creating order")
	}

	return nil
}
