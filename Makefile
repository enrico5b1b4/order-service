.PHONY: imports sec build run test integration-test docker_build docker_run testdb-up testdb-down

PORT := 8001
COMPLETE_ORDER_CALLBACK_URL := "http://localhost:8001/complete_order"
POSTGRES_DSN := "postgres://postgres:example@localhost:5432/orderservice?sslmode=disable"
ORDER_PROCESS_SERVICE_URL := "http://localhost:8000"
POSTGRES_TEST_DSN := "postgres://postgres:example@localhost:5433/orderservice?sslmode=disable"

sec:
	@gosec ./...

imports:
	@goimports -w .

build: imports
	@go build

run: build
	@PORT=$(PORT) \
	POSTGRES_DSN=$(POSTGRES_DSN) \
	ORDER_PROCESS_SERVICE_URL=$(ORDER_PROCESS_SERVICE_URL) \
	COMPLETE_ORDER_CALLBACK_URL=$(COMPLETE_ORDER_CALLBACK_URL) \
	./order-service

docker_build: build
	docker image rm -f enrico5b1b4/orderservice_app
	docker build -t enrico5b1b4/orderservice_app .

docker_run:
	docker run \
	--net="host" \
	-e PORT=$(PORT) \
	-e COMPLETE_ORDER_CALLBACK_URL=$(COMPLETE_ORDER_CALLBACK_URL) \
	-e POSTGRES_DSN=$(POSTGRES_DSN) \
	-e ORDER_PROCESS_SERVICE_URL=$(ORDER_PROCESS_SERVICE_URL) \
	enrico5b1b4/orderservice_app:latest

test:
	@go test ./...

integration-test:
	@POSTGRES_TEST_DSN=$(POSTGRES_TEST_DSN) go test -count=1 -p 1 ./...

testdb-up:
	@docker stop orderservice_test_postgres96 || true && docker rm  orderservice_test_postgres96 || true
	@docker run --name orderservice_test_postgres96 \
				-p 5433:5432 \
				-e POSTGRES_PASSWORD=example \
				-e POSTGRES_USER=postgres \
				-e POSTGRES_DB=orderservice \
				-d postgres:9.6
	@sleep 5

testdb-down:
	@docker rm -f orderservice_test_postgres96
