include .env

create-container:
	docker run --name ${DB_CONTAINER_NAME} -p ${DB_PORT}:${DB_PORT} -e POSTGRES_USER=${POSTGRES_USER} \
	 -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} -d postgres:12-alpine

create-db:
	docker exec -it ${DB_CONTAINER_NAME} createdb --username=${POSTGRES_USER} --owner=${POSTGRES_USER} ${DB_NAME}

# generate graphql based on schema.graphql
graphql-generate:
	@echo "Generating graphql go files..."
	@go run github.com/99designs/gqlgen generate
	@echo "Generate successful"

# add-migration named 'init'
add-migration-init:
	@sqlx migrate add -r init

# applies all migrations
migrate-up:
	@sqlx migrate run --database-url "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

# rollback one migration
migrate-down:
	@sqlx migrate revert --database-url "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

build:
	@if [ -f "${BINARY_NAME}" ]; then \
		rm ${BINARY_NAME}; \
	fi
	@echo "Building server binary..."
	@go build -o ${BINARY_NAME} ./server/*.go
	
run: build 
	@echo "Starting server..."
	@./${BINARY_NAME}

stop:
	@echo "Stopping server..."
	@-pkill -SIGTERM -f "./${BINARY_NAME}"
	@echo "Server stopped..."