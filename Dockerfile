# build stage
FROM golang:1.12.5 as builder

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .

# final stage
FROM scratch
COPY --from=builder /app/order-service /app
COPY --from=builder /app/migrations /migrations/
EXPOSE 8001
ENTRYPOINT ["/app"]