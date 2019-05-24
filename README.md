# Order Service

Orders service which interacts with a pre-existing order process service.

## Getting Started

### Prerequisites
To run and test the project you will need the following to be installed

- docker
- docker-compose
- go
- goimports (`go get golang.org/x/tools/cmd/goimports`)
- gosec (`go get github.com/securego/gosec/cmd/gosec`)


### Installing

Build the image `enrico5b1b4/orderservice_app`

```
make docker_build
```
or pull from dockerhub

```
docker pull enrico5b1b4/orderservice_app
```
Run all the services with docker-compose

```
docker-compose up
```
The order service will be running on `localhost:8001`
## Running the tests

First run `go generate` to generate all the required mocks

```
go generate
```

### Unit tests
Run the unit tests with
```
make test
```
### Integration tests
Integration tests rely on a test postgres database to be up and running.
To start the test postgres database
```
make testdb-up
```
Then run the integration tests
```
make integration-test
```

# API endpoints

These endpoints allow you to create and query orders.  
Please note: order ids must be uuids ([Generate UUIDs online here](https://www.uuidgenerator.net/version4))

## GET  
[/order](#get-order)  
[/order/:orderId](#get-orderorderid)  


## POST  
[/order](#post-order)  
[/process_order](#post-process_order)  
[/complete_order](#post-complete_order)  

___
### [GET] /order
Returns all orders that meet the supplied filter(s). If no filter parameters are supplied then all orders should be returned.

**Parameters**

|          Name | Required |  Type   | Description                                                                                                                                                           |
| -------------:|:--------:|:-------:| --------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
|     `status`  |          | string  | filters orders by status (`CREATED`, `PROCESSING`, `FULFILLED`, `FAILED`)                                                                     |

**Response**

```
{
    "orders": [
        {
            "id": "a970d7c6-057b-42e7-8b15-3dbe3fecca17",
            "status": "PROCESSING"
        },
        {
            "id": "c6f1b778-711e-4ea7-b1f4-83e8400bddbe",
            "status": "FAILED"
        },
        {
            "id": "0867c280-b8cd-4020-8caa-7b4af6f3dfd5",
            "status": "PROCESSING"
        },
        {
            "id": "63efb280-205f-4e36-a581-61ace8e9212a",
            "status": "CREATED"
        },
        {
            "id": "6e87d902-1e46-4cff-bd54-ccfb8244d4bc",
            "status": "FULFILLED"
        },
        {
            "id": "a593e5b2-0c91-4b96-8b30-1ef3d0236ec6",
            "status": "FAILED"
        }
    ]
}
```
___
### [GET] /order/:orderId
Returns an order by id.

**Response**

```
{
    "order": {
        "id": "c6f1b778-711e-4ea7-b1f4-83e8400bddbe",
        "status": "FAILED"
    }
}
```
___

### [POST] /order
Creates an order

**Body**

|          Name | Required |  Type   | Description                                                                                                                                                           |
| -------------:|:--------:|:-------:| --------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
|     `id` | required | string  | Id of the order to create                                                                     |

**Request**

```
{
	"id": "a593e5b2-0c91-4b96-8b30-1ef3d0236ec6"
}
```
___

### [POST] /process_order
Starts processing an order previously created

**Body**

|          Name | Required |  Type   | Description                                                                                                                                                           |
| -------------:|:--------:|:-------:| --------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
|     `id` | required | string  | Id of the order to process                                                                     |

**Request**

```
{
	"id": "a593e5b2-0c91-4b96-8b30-1ef3d0236ec6"
}
```
___

### [POST] /complete_order
Completes an order

**Body**

|          Name | Required |  Type   | Description                                                                                                                                                           |
| -------------:|:--------:|:-------:| --------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
|     `order_id`| required | string  | Id of the order to complete                                                                     |
|     `status`  | required | string  | Status of the order process of the order (`SUCCEEDED`, `FAILED`)                                                                |

**Request**

```
{
	"order_id": "a970d7c6-057b-42e7-8b15-3dbe3fecca17",
	"status": "SUCCEEDED"
}
```