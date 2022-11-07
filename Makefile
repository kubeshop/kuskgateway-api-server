.PHONY: all
all: format test build

.PHONY: build
build:
	docker buildx build -t kusk-gateway-api -f build/api-server/Dockerfile .
	docker buildx build -t kusk-gateway-api-websocket -f build/websocket/Dockerfile .

server-generate:
	openapi-generator-cli generate -i api/openapi.yaml -g go-server -o server/ --additional-properties=featureCORS=true

run:
	docker-compose up --build -d

run-minikube:
	docker-compose -f docker-compose.yaml -f docker-compose-minikube.yaml up --build --force-recreate

test:
	cd ./server && FAKE=true go test -v -count=1 ./...

.PHONY: format
format:
	cd ./server && go fmt ./...
