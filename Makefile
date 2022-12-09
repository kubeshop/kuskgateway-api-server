FAKE ?= true

.PHONY: all
all: format test build

.PHONY: build
build:
	docker buildx build -t kubeshop/kusk-gateway-api -f build/api-server/Dockerfile .
	docker tag kubeshop/kusk-gateway-api ttl.sh/kubeshop/kusk-gateway-api:latest
	docker tag kubeshop/kusk-gateway-api ttl.sh/kubeshop/kusk-gateway-api:$(shell git describe --tags $(shell git rev-list --tags --max-count=1))
	@echo
	docker buildx build -t kubeshop/kusk-gateway-api-websocket -f build/websocket/Dockerfile .
	docker tag kubeshop/kusk-gateway-api-websocket ttl.sh/kubeshop/kusk-gateway-api-websocket:latest
	docker tag kubeshop/kusk-gateway-api-websocket ttl.sh/kubeshop/kusk-gateway-api-websocket:$(shell git describe --tags $(shell git rev-list --tags --max-count=1))

server-generate:
	openapi-generator-cli generate -i api/openapi.yaml -g go-server -o server/ --additional-properties=featureCORS=true

run:
	docker-compose up --build -d

run-minikube:
	docker-compose -f docker-compose.yaml -f docker-compose-minikube.yaml up --build --force-recreate

test:
	cd ./server && FAKE=${FAKE} go test -v -count=1 ./...

.PHONY: format
format:
	cd ./server && go fmt ./...
