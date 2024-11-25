include .dev.env

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
	@migrate create -ext sql -dir database/migrations -seq init

# applies all migrations
migrate-up:
	@migrate -database ${POSTGRESQL_URL} -path database/migrations up

# rollback one migration
migrate-down:
	@migrate -database ${POSTGRESQL_URL} -path database/migrations down 1

test:
	@echo ${DBSTRING}

# build and run go with Air
run:
	@echo "Starting server..."
	@air c .air.toml

# use remote docker context
context-remote:
	@docker context use do-droplet 
	eval "$(ssh-agent -s)"
	@ssh-add ~/.ssh/id_rsa
	@docker info 

# for rebuildling docker compose
docker-rebuild:
	@docker-compose down
	@docker-compose up --build

docker-prune:
	@docker image prune -f
	@docker volume prune -f
	@docker container prune -f