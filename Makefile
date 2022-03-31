build:
	docker build --network=host -t kusk-gateway-api ./server

server-generate:
	openapi-generator-cli generate -i api/openapi.yaml -g go-server -o server/ --additional-properties=featureCORS=true
run:
	docker-compose up --build -d 