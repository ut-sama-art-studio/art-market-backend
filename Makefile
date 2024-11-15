include .env

create-db-container:
	docker run --name ${DB_CONTAINER_NAME} -p ${DB_PORT}:${DB_PORT} -e POSTGRES_USER=${POSTGRES_USER} \
	 -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} -d postgres:12-alpine

create-db:
	docker exec -it ${DB_CONTAINER_NAME} createdb --username=${POSTGRES_USER} --owner=${POSTGRES_USER} ${DB_NAME}

# generate graphql based on schema.graphql
graphql-generate:
	go get github.com/99designs/gqlgen@v0.17.47
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

test:
	@echo ${DBSTRING}

run:
	@echo "Starting server..."
	@air c .air.toml

check:
	@echo $(shell pwd)

stop:
	@docker stop ${API_DOCKER_NAME}