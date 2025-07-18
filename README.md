
## Backend development Linux/WSL setup
Install [Golang 1.23](https://go.dev/dl/)

Install the required linux packages:
```bash
sudo apt install make build-essential libssl-dev pkg-config
```

Install [air](https://github.com/air-verse/air) for hot reload && [migrate](https://github.com/golang-migrate/migrate) for migrations:
```bash
go install github.com/air-verse/air@latest
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

To view Postgres DB, install the [Postgres extension](https://marketplace.visualstudio.com/items?itemName=ckolkman.vscode-postgres) in VSC, or a SQL manager like DBeaver

## SQL Docker container setup
Install [Docker Engine](https://docs.docker.com/engine/install/) & [Docker Desktop](https://docs.docker.com/desktop/setup/install/linux/)

Get the secrets files (`.dev.env`, `.env`) from an admin and place them in the art-market-backend directory.

Run:

Start docker service (if not already running)
```bash
systemctl start docker
```

```bash
make create-db-container

make create-db
```

## Running the backend:
```bash
make run
```

## GraphQL development
The project uses 99designs/gqlgen package to handle generating code needed for GraphQL

To add new api:
1. Change the GraphQL schemas, in '/graph/schemas/'
2. Run the following command to generate the corresponding resolvers and models
```bash 
make graphql-generate
```
3. Implementing the resolver function
4. Test with GraphQL playground at http://localhost:8080/api

## Deployment 
Connect to remote Digital Ocean droplet with Docker 
```bash
eval "$(ssh-agent -s)"
make context-remote # prompts password input
```

Build Docker image
```bash
make docker-rebuild
```

(Optional) swich back to default context
```
docker context use default
```


