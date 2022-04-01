FROM golang:1.17 AS build
WORKDIR /go/src
COPY . .
COPY ./kubeconfig .
ENV CGO_ENABLED=0
RUN go get -d -v ./...

RUN go build -a -installsuffix cgo -o kusk-gateway-api .

FROM scratch AS runtime
COPY --from=build /go/src/kusk-gateway-api ./
COPY --from=build /go/src/kubeconfig ./
ENV KUBECONFIG=kubeconfig
EXPOSE 8080/tcp
ENTRYPOINT ["./kusk-gateway-api"]
