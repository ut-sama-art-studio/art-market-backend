include .dev.env

create-db-container:
	docker run -it --rm --add-host=host.docker.internal:host-gateway --name ${DB_CONTAINER_NAME} -p ${DB_PORT}:${DB_PORT} -e POSTGRES_USER=${POSTGRES_USER} \
	 -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} -d postgres:12-alpine

create-db:
	docker exec -it ${DB_CONTAINER_NAME} createdb --username=${POSTGRES_USER} --owner=${POSTGRES_USER} ${DB_NAME}

# generate graphql based on schema.graphql
graphql-generate:
	go get github.com/99designs/gqlgen@v0.17.47
	@echo "Generating graphql go files..."
	@go run github.com/99designs/gqlgen generate
	@echo "Generate successful"

# add-migration given NAME=" "
add-migration:
	@migrate create -ext sql -dir database/migrations -format unix $(NAME)   

# applies all migrations
migrate-up:
	@migrate -database ${POSTGRESQL_URL} -path database/migrations up

# rollback one migration
migrate-down:
	@migrate -database ${POSTGRESQL_URL} -path database/migrations down 1

# check current migration version
migrate-version:
	@migrate -database ${POSTGRESQL_URL} -path database/migrations version

# force apply migration given VERSION
migrate-force:
	@migrate -database ${POSTGRESQL_URL} -path database/migrations force $(VERSION)   

#get the POSTGRESQL_URL
URL:
	@echo ${POSTGRESQL_URL}


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

# rebuilds with docker-compose
docker-rebuild:
	@docker-compose down
	@docker-compose up --build -d

# removes all docker files
docker-prune:
	@docker image prune -f
	@docker volume prune -f
	@docker container prune -f

# connects to Digital Ocean droplet using private key
ssh:
	@ssh -i ~/.ssh/dig_ocean_art_market_ssh root@143.198.44.171