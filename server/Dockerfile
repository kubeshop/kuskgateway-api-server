FROM --platform=$BUILDPLATFORM docker.io/golang:1.17 as builder
WORKDIR /go/src
# Copy `go.mod` for definitions and `go.sum` to invalidate the next layer
# in case of a change in the dependencies
COPY go.mod go.sum ./
# Download dependencies
RUN go mod download

ARG TELEMETRY_TOKEN
ARG TARGETARCH
ARG TARGETOS

COPY . .
RUN RUN GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 go build -v -ldflags "-X github.com/kubeshop/kusk-gateway/pkg/analytics.TelemetryToken=$TELEMETRY_TOKEN" -o kusk-gateway-api

FROM --platform=$BUILDPLATFORM gcr.io/distroless/static:nonroot
COPY --from=build --chown=65532:65532 /go/src/kusk-gateway-api ./
EXPOSE 8080/tcp

USER 65532:65532

ENTRYPOINT ["./kusk-gateway-api"]
