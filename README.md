
## Backend development Linux/WSL setup
Have Golang installed

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

## Running 
```bash
make run
```

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


