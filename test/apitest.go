package test

import (
	"github.com/enrico5b1b4/order-service/app"
	"github.com/steinfletcher/apitest"
)

func ApiTest() *apitest.APITest {
	dbRaw := DBConnect()
	orderServiceApp := app.New(
		dbRaw,
		"http://orderprocessservice.local",
		"http://localhost/complete_order",
	)

	return apitest.New("").
		Handler(orderServiceApp)
}
