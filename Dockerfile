FROM golang:1.24.5-alpine3.22 AS build

WORKDIR /go/src/app

ENV PORT=8080
ENV GOCACHE=/root/.cache/go-build

# Copy config files
COPY ./config ./config

# Check for any changes to dependencies.
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum

# Download dependencies
RUN --mount=type=cache,target=/go/pkg/mod go mod download

# Copy main
COPY ./main.go ./main.go

# Copy migrations
COPY ./migrations ./migrations

# Copy source code
COPY ./internal ./internal

RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 GOOS=linux go build -o /go/bin/app

FROM gcr.io/distroless/static-debian11
COPY --from=build /go/src/app /
COPY --from=build /go/bin/app /
ENV ENV=QA

CMD ["/app"]