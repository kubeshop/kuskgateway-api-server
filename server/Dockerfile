FROM golang:1.17 AS build
WORKDIR /go/src
# Copy `go.mod` for definitions and `go.sum` to invalidate the next layer
# in case of a change in the dependencies
COPY go.mod go.sum ./
# Download dependencies
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v -o kusk-gateway-api

FROM scratch AS runtime
COPY --from=build --chown=65532:65532 /go/src/kusk-gateway-api ./
EXPOSE 8080/tcp

USER 65532:65532

ENTRYPOINT ["./kusk-gateway-api"]
