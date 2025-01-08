FROM golang:1.23.2-alpine3.20 AS tester

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
ENTRYPOINT ["go", "test", "-coverprofile=coverage/coverage.html", "./..."]
# CMD ["-v"]
