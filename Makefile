build:
	@go build -o bin/app users-microservice/cmd/server

run:
	@docker-compose -f docker-compose.yaml up

test-db-up:
	@docker-compose -f docker-compose.test.yaml up -d
	echo "Waiting for database to be ready..."
	@sleep 5

test-db-down:
	@docker-compose -f docker-compose.test.yaml down

test-integration: test-db-up
	echo "Running integration tests..."
	@go test ./tests/integration/... -v
	@make test-db-down

clean:
	rm -rf bin/
	@make test-db-down