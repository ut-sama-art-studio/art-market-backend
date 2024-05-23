
## Backend development Linux/WSL setup
Have Golang installed

Install required packages
```
sudo apt install make  
sudo apt install build-essential
sudo apt install libssl-dev 
sudo apt install pkg-config
```

Install Rust
```
curl https://sh.rustup.rs -sSf | sh
```

Install `sqlx-cli` with Rust to manage migrations
```
cargo install sqlx-cli 
```

To view Postgres DB, install the [Postgres extension](https://marketplace.visualstudio.com/items?itemName=ckolkman.vscode-postgres) in VSC, or a SQL manager like DBeaver