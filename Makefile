build:
	docker build -t kusk-gateway-api server

server-generate:
	openapi-generator-cli generate -i api/openapi.yaml -g go-server -o server/ --additional-properties=featureCORS=true

run:
	docker-compose up --build -d

run-minikube:
	docker-compose -f docker-compose.yaml -f docker-compose-minikube.yaml up --build --force-recreate

test:
	cd ./server && FAKE=true go test -v ./...
