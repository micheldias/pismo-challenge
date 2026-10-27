# Pismo Challenge API

A RESTful API built with Go and PostgreSQL to manage accounts and financial transactions.

The application allows users to create accounts, retrieve account information, and register transactions associated with an account.

---
## Technologies & Tools

___

[![](https://img.shields.io/badge/Code-Go-informational?style=plastic&logo=Go&logoColor=white&color=2bbc8a)](https://golang.org/)
[![](https://img.shields.io/badge/PostgreSQL-316192?logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![](https://img.shields.io/badge/Tools-Docker-informational?style=plastic&logo=Docker&logoColor=white&color=yellow)](https://www.docker.com/)
[![](https://img.shields.io/badge/docker-compose-blue?logo=docker)](https://docs.docker.com/compose/)


## Architectural Overview.
___
The project was structured with simplicity and loose coupling in mind to facilitate maintenance, test creation, and future expansion.

## Architectural Decisions
___
- **Package separation by responsibility**: The layers are designed to be well separated and loosely coupled.
- **Centralized Error Handling**: Error logic is handled in a single flow, ensuring greater predictability in API behavior.
- **Automated Testing**: Robust tests were implemented for the application.
- **Docker**: Docker is used to initialize the application. There is no need to install any project dependencies on your computer.

### Package structure
___
```plaintext
pismo-challenge/
│
├── cmd/
│   ├── api-server/            # Application entry point for API server.
├── docs/                      # API documentation (e.g., Swagger files).
├── internal/
│   ├── controller/            # API layer response for handle web request/response (controllers, request/response handling).
│   ├── infra/                 # Infrastructure layer (technical concerns and cross-cutting components).
│   │   ├── context/           # context handler helper functions.
│   │   └── database/          # Database configuration related.
│   │   ├── logging/           # Configuration for structred log.
│   │   └── server/            # HTTP server setup, routes, helper middlewares and server-related models.
│   │   ├── settings/          # env vars.
│   ├── model/                 # Application domain models
│   ├── repository/            # Data persistence layer (interfaces and implementations for accessing DB)
│   └── service/               # Business logic and main use case implementations.
└── mocks/                     # Mocks and startup scripts for the local development environment.  
└── README.md                  # Project documentation
└── run.md                     # Setup instrunctions
```

## Endpoints

A swagger was made available to better test the API, as well as visualize the possible responses and request and response contracts.
You can access by :
http://localhost:8080/swagger/index.html

You must start up the app. You can get instructions in [run.md](run.md)