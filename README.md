# URL shortener & resolver Web API app

## Pre-requisites
- Go v1.16+
- [Echo Web Framework](https://echo.labstack.com)
- Visual Studio Code with Go extension installed globally
- PostgreSQL v13+
- Docker v20.10.6 build 370c289 or above
- Docker Compose v3.7+ 

### Health Check Endpoint
```
http://localhost:3355/
```
### Swagger API Documentation Available At
```
http://localhost:3355/swagger/index.html
```
## API Endpoints
| Name  | Path | Method |
|-------|------|--------|
| Generator | /api/generate | **POST** |
| Resolver | /api/resolve | **GET** |

## Environment Variables
| Name  | Meaning |
|-------|---------|
| LOG_MINIMUM_LEVEL | Minimum logging level |
| CUSTOM_PORT_NUMBER | Server port binding |
| USE_IN_MEMORY_DB | [Y/N] Y to use in-memory database, N to use PostgreSQL |
| POSTGRES_USER | PostgreSQL database server user name |
| POSTGRES_PASSWORD | PostgreSQL database server password |
| POSTGRES_HOST | PostgreSQL database server name |
| POSTGRES_PORT | PostgreSQL database server port |
| POSTGRES_DATABASE | PostgreSQL database name |
| POSTGRES_SSL_MODE | PostgreSQL SSL mode info, keep blank if SSL is not enabled |

## Features
- High performance, fast APIs with minimal resource foot print
- Data persistence in PostgreSQL
- PostgreSQL database migrations
- Open API documentation support using Swagger
- GitHub actions

## Open API Swagger Documentation Generation / Updation
First, add / edit inline documentation comments following [Declarative Comments Format](https://github.com/swaggo/swag#declarative-comments-format). After this follow:
```powershell
cd go-url-shortener\src
swag init
```
This will generate/re-generate documentation artifacts inside `go-url-shortener\src\docs` directory.

## Performing Database Migration
### Creating New Database Migration
```powershell
cd go-url-shortener\src
go run Database\Migrations\main.go create <NAME_OF_MIGRATION>
```

After this a new `Go` file will be generated at `Database\Migrations` directory with the name: `<TIMESTAMP>_<NAME_OF_MIGRATION>`, Update the created file with `up` & `down` scripts.

### Running Database Migration
Update the environment variables of .env file as per configuration requirement, select the configuration named: **Database Migration**, and start running the application in VS Code.

## TODOs
- HTTPS support / redirection