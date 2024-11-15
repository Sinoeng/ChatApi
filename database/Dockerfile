# builder
FROM golang:1.23.2-alpine3.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ChatApi

# deployer
FROM alpine:3.20.3 AS deployer

WORKDIR /app

COPY --from=builder /app/program ./program
COPY .env ./.env

ENTRYPOINT ["./ChatApi"]
