# Advanced Golang Api

Hello! This project is meticulously crafted using the Go programming language and incorporates best practices for robust and efficient code.

## Features

- Implements advanced algorithms for efficient processing.
- Incorporates clean and maintainable code architecture for easy scalability.

## Techs

- Docker
- Redis
- PostgreSQL
- Golang
- Air (hot reload)
- Migrate

## Getting Started

To run this application, you need Docker installed on your system.

1. Clone the repository to your local machine:

```sh
git clone https://github.com/henriqueassiss/advanced-golang-api.git
```

2. Navigate to the project directory:

```sh
cd advanced-golang-api
```

3. Create a .env file like .env.example, but with it's values, here is an example:

```sh
# App
APP_ENVIRONMENT=development

# Api
API_HOST=0.0.0.0
API_PORT=8080
API_SECRET=some_long_string
API_READ_HEADER_TIMEOUT=60s
API_GRACEFUL_TIMEOUT=8s
API_REQUEST_LOG=true

# Client
CLIENT_BASE_URL=http://localhost:3000

# Cors
CORS_ALLOWED_ORIGINS=http://localhost:3000

# Database
DB_DRIVER=pgx
DB_HOST=db
DB_PORT=5432
DB_NAME=postgres
DB_USER=postgres
DB_PASSWORD=123456Abc
DB_SSL_MODE=disable
DB_MAX_CONNECTION_POOL=85
DB_MAX_IDLE_CONNECTIONS=85
DB_CONNECTIONS_MAX_LIFE_TIME=300s

# Redis
CACHE_ADDRESS=localhost:6379
CACHE_PORT=6379
CACHE_PASSWORD=
CACHE_DB=0
```

4. Build and start the application using Docker Compose:

```sh
docker-compose up
```

5. Access the application:
   Once the Docker containers are up and running, you can access the application at http://localhost:your_port.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

Hope you like it!
