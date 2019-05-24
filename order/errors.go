package order

import "github.com/enrico5b1b4/order-service/errors"

var OrderNotFoundError = &errors.Error{
	Code:    "3",
	Message: "Order not found.",
}

var InvalidOrdersStatusError = &errors.Error{
	Code:    "4",
	Message: "Invalid status filter supplied. Must be one of CREATED, PROCESSING, FULFILLED or FAILED",
}

var OrderAlreadyExistsError = &errors.Error{
	Code:    "5",
	Message: "Order id already exists.",
}

var UnknownProcessOrderStatusError = &errors.Error{
	Code:    "6",
	Message: "Unknown process order status.",
}

var OrderAlreadyBeingProcessedError = &errors.Error{
	Code:    "7",
	Message: "Cannot start processing order since it is already being processed.",
}

var OrderAlreadyCompleteError = &errors.Error{
	Code:    "8",
	Message: "Cannot process a completed order.",
}
