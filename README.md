## Bank API

This project is a REST API for a banking system, built with Go and the Gin-Gonic framework. It uses Docker for containerization and PostgreSQL for data persistence. The API supports user authentication and authorization using JWT.

## Features
#### User Authentication and Authorization
Users can sign up, log in, and perform various operations based on their access rights.
#### Bank Account Management
Users can create, read, update, and delete bank accounts with different currencies.
#### Transaction History
Users can view their transaction history, including deposits, withdrawals, and transfers.
#### Transfers
Users can transfer funds between different accounts.

## Running locally
Start the Docker services:

```docker compose up```

## Docs

Endpoints and their description can be seen in swagger ui:
```
http://localhost:8080/swagger-ui/index.html/
```

## Contributing

Contributions are what make the open-source community such an amazing place to learn, inspire, and create. Any contributions you make are greatly appreciated.