FROM golang:1.23.3-alpine3.20 AS builder

WORKDIR /app

# Check for any changes to dependencies.
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum

# Download dependencies
RUN --mount=type=cache,target=/go/pkg/mod go mod download

# Copy main
COPY ./main.go ./main.go

RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 GOOS=linux go build -o main .

CMD ["/app/main"]