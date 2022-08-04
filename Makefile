BINARY=service

# Actions to execute locally
build:
	go build -o ${BINARY} cmd/mocks/*.go

generate-swagger:
	swag init -g cmd/mocks/main.go

generate-mocks:
	go generate ./...

test:
	go test -cover -race ./...

generate-web:
	./scripts/build-web.sh

run-mysql:
	docker-compose -f ./scripts/local/mocks-mysql/docker-compose.yml up -d --remove-orphans --build
	sleep 15

run-file:
	docker-compose -f ./scripts/local/mocks-file/docker-compose.yml up -d --remove-orphans --build
	docker ps

stop-mysql:
	docker-compose -f ./scripts/ingestor-services/docker-compose.yml stop
	docker ps

stop-file:
	docker-compose -f ./scripts/ingestor-services/docker-compose.yml stop
	docker ps

pre-commit:
	./scripts/pre-commit.sh
