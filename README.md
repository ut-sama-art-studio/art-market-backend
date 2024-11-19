
## Backend development Linux/WSL setup
Have Golang installed

Install required linux packages
```
sudo apt install make  
sudo apt install build-essential
sudo apt install libssl-dev 
sudo apt install pkg-config
```

Install Air for hot reload && golang-migrate for migrations
```
go install github.com/air-verse/air@latest
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

To view Postgres DB, install the [Postgres extension](https://marketplace.visualstudio.com/items?itemName=ckolkman.vscode-postgres) in VSC, or a SQL manager like DBeaver

## Running 
```
make run
```