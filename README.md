
## Backend development Linux/WSL setup
Install Golang 1.23 

Install required linux packages
```bash
sudo apt install make  
sudo apt install build-essential
sudo apt install libssl-dev 
sudo apt install pkg-config
```

Install Air for hot reload && golang-migrate for migrations
```bash
go install github.com/air-verse/air@latest
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

To view Postgres DB, install the [Postgres extension](https://marketplace.visualstudio.com/items?itemName=ckolkman.vscode-postgres) in VSC, or a SQL manager like DBeaver

## SQL Docker container setup
Install Docker & Docker Desktop

run 
```bash
make create-db-container

make create-db
```

## Running 
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


