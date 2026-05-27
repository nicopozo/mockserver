BINARY=service
VERSION?=$(shell cat VERSION)

# Actions to execute locally
build:
	go build -o ${BINARY} cmd/mocks/*.go

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
	docker-compose -f ./scripts/local/mocks-json/docker-compose.yml up -d --remove-orphans --build
	docker ps

stop-mysql:
	docker-compose -f ./scripts/local/mocks-mysql/docker-compose.yml stop
	docker ps

stop-file:
	docker-compose -f ./scripts/local/mocks-json/docker-compose.yml stop
	docker ps

pre-commit:
	./scripts/pre-commit.sh

# Docker Image Management
docker-build:
	chmod +x ./scripts/docker-build.sh
	./scripts/docker-build.sh $(VERSION)

docker-hub-push:
	chmod +x ./scripts/docker-hub-push.sh
	./scripts/docker-hub-push.sh $(VERSION)

# AWS Infrastructure
aws-create-role:
	chmod +x ./scripts/aws/create-lambda-role.sh
	./scripts/aws/create-lambda-role.sh

aws-init-db:
	chmod +x ./scripts/aws/create-dynamo-tables.sh
	./scripts/aws/create-dynamo-tables.sh

aws-enable-api-gateway:
	chmod +x ./scripts/aws/enable-api-gateway.sh
	./scripts/aws/enable-api-gateway.sh

aws-docker-push:
	chmod +x ./scripts/aws/docker-push.sh
	./scripts/aws/docker-push.sh $(VERSION)

# AWS Lambda Deployment
aws-lambda-deploy:
	chmod +x ./scripts/aws/deploy-lambda.sh
	./scripts/aws/deploy-lambda.sh $(VERSION)

aws-lambda-full: bump docker-build aws-docker-push aws-lambda-deploy

bump:
	@old_version=$$(cat VERSION); \
	major=$$(echo $$old_version | cut -d. -f1); \
	minor=$$(echo $$old_version | cut -d. -f2); \
	patch=$$(echo $$old_version | cut -d. -f3); \
	if [ "$(PART)" = "major" ]; then \
		new_version=$$((major+1)).0.0; \
	elif [ "$(PART)" = "minor" ]; then \
		new_version=$$major.$$((minor+1)).0; \
	else \
		new_version=$$major.$$minor.$$((patch+1)); \
	fi; \
	echo $$new_version > VERSION; \
	echo "🚀 Bumped version from $$old_version to $$new_version (part: $${PART:-patch})"
