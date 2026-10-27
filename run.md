## Project Setup

You can choose to run this project via docker compose, docker or locally.

### Swagger

A swagger was made available to better test the API, as well as visualize the possible responses and request and response contracts.
You can access by :

http://localhost:8080/swagger/index.html

If you don't have swaggo installed, you can use the command

```
go install github.com/swaggo/swag/cmd/swag@latest
```

If you want to regenerate the documentation, simply run the following command in the project root folder,

```
swag init -d ./cmd/api-server,./internal
```
To format the Swagger notes, run the following command
```
swag fmt -d ./internal,./cmd/api-server
```

### Running the application with Docker Compose

Docker Compose allows you to run the application and its database together without installing Go or PostgreSQL locally.

#### Prerequisites

Make sure you have the following tools installed:

- [Docker Desktop](https://www.docker.com/products/docker-desktop/) or Docker Engine.
- Docker Compose v2 or later.

#### Starting the application

From the project root directory, run the following command to build the application image and start the application and database containers:

```bash
docker compose up --build
```

To run the containers in the background, use:

```bash
docker compose up --build -d
```

#### Checking if the application is running

Once the containers are up, check the application's health endpoint:

```bash
curl http://localhost:8080/health
```

#### Viewing application logs

To follow the logs of the running containers, execute:

```bash
docker compose logs -f
```

To view logs for the application service only, use:

```bash
docker compose logs -f web-app
```

Replace `web-app` with the application service name defined in your `docker-compose.yml`, if different.

#### Stopping the application

To stop and remove the containers, run:

```bash
docker compose down
```

This command removes the containers and network created by Docker Compose. Database data may persist if a named volume is configured.

To stop the application and remove its containers, network, and database volumes, run:

```bash
docker compose down -v
```

**Warning:** The `-v` option removes the database volumes and may permanently delete persisted database data.


### Running the Docker application

In order to run the docker application, you must have the following tools configured:
  After that, you must be good to go!

## Building the Image
___
In the root folder,you must run the command:

```bash
docker build -t pismo-challeng-api:1.0 .
```

## Running the Image
To start the application, execute the command as follows:

```
docker run -p 8080:8080 pismo-challeng-api:1.0

```
Ensure you have a database installed and configured. You will need to adjust the database credentials using environment variables before running the application

If you want to change any application variable, you must run the start command with the following options:

```
docker run -p 8080:8080 -e APP_APPNAME=new_app_name pismo-challeng-api:1.0
```

### Checking if the application is UP

You can check if the application is running by calling its health endpoint:

```
curl --location --request GET 'http://localhost:8080/health'
```

### Running the application locally

In order to run the application, in your local machine, you must have the following tools configured:

- Go 1.27.3 or later ([installation guide](https://go.dev/doc/install))

After that, you must be good to go!

### Downloading Dependencies
In the project root folder, run the command below to install the dependencies

```bash
go mod tidy
```

In the root folder, to run the server use:
```bash
go run cmd/api-server/main.go
```

In the root folder use, if you want to run the tests use:
```bash
go test ./...
```