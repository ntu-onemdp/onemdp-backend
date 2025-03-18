# onemdp-backend

This repository hosts the backend for OneMDP. It serves as the connector between the frontend and the Postgres database.

## Requirements

- [Docker](https://www.docker.com/): The easiest way to run the whole application is with Docker compose from the main repository.
- [Golang 1.23.3](https://go.dev/doc/install) (Optional if you are running from Docker)
- [Postgres](https://www.postgresql.org/) instance running on port 5432. (Optional if using Docker compose)

## Getting started

### 1. Set up Postgres instance

The easiest way to set up the Postfres instance is running the docker compose command in the [parent repository](https://github.com/ntu-onemdp/onemdp). It is highly recommended to create a separate postgres user instead of using the default postgres user.

Ensure that you have the following:

1. Your postgres username and password
2. A new database (if you are using the values from the template file, it will be `dev_2`)

### 2. Set up `.env` file

Remove `.template` from the filename of `.env.template` in the config folder. Set the values of `PG_USERNAME` and `PG_PASSWORD` to the values you have used when setting up postgres.

Set `POSTGRES_NETLOC` to where you are running the backend from. The default value assumes you are running the backend as a standalone container.

### 3. Set up JWT key

#### Generate JWT secret key

Generate your JWT secret key and copy the generated key into `config/jwt-key.txt`. The most secure way to generate the secret key is with the following command:

```sh
# Note: You will need to clean up the .txt file after generation
openssl genrsa -out ./config/jwt-key.txt 4096
```

Alternatively, use an [online key generator](https://jwtsecret.com/generate) (NOT recommended for production!).

> [!WARNING]
> The JWT key and database password should be kept secret and never shared with anyone.

## Running the code

There are 3 ways to run the backend:

1. Locally (requires go to be installed)
2. As a standalone docker container
3. With docker compose from the [parent repository](https://github.com/ntu-onemdp/onemdp)

Before running, ensure that `PSTGRES_NETLOC` is set to the correct value in `.env`.

### Running locally

1. Download dependencies

```sh
go mod download
```

2. Run

```sh
go run main.go
```

### Running as a standalone docker container

Simply run:

```sh
docker rm -f onemdp-dev-1 || true && DOCKER_BUILDKIT=1 docker build -f Dockerfile.dev  -t onemdp-dev-1 . && docker run -it -p 8080:8080 --name onemdp-dev-1 onemdp-dev-1
```

OR

```sh
sh docker.sh
```
