## Bank API 💰

This project is a REST API for a banking system, built with Go and Gin. It uses JWT for
authentication and PostgreSQL for data persistence.

## Features 

Users can sign up, login, and perform various operations based on their access rights, including creation of bank
accounts with different currencies, money deposits/withdraws/transfers. All the transactions can be listed in history.

Some mock users and their bank accounts are added by default [here](https://github.com/DimaGitHahahab/bank-api/tree/main/migrations). 

## Running locally 

1. Edit ```.env``` file for configuration (optional)
2. Start the Docker services: ```docker compose up```

## Docs 

Endpoints and their description can be seen in Swagger ui:

```
http://localhost:8080/swagger-ui/index.html/
```

## Contributing 💍

Contributions are what make the open-source community such an amazing place to learn, inspire, and create. Any
contributions you make are greatly appreciated.
