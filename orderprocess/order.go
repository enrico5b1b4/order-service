package orderprocess

import uuid "github.com/satori/go.uuid"

type ProcessOrderStatus = string

const (
	RUNNING   ProcessOrderStatus = "RUNNING"
	SUCCEEDED ProcessOrderStatus = "SUCCEEDED"
	FAILED    ProcessOrderStatus = "FAILED"
)

type Order struct {
	ID     uuid.UUID          `json:"order_id"`
	Status ProcessOrderStatus `json:"status,omitempty"`
}
