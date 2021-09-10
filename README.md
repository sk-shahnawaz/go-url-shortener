# URL shortener & resolver Web API app

## Pre-requisites
- Go v1.16+
- [Echo Web Framework](https://echo.labstack.com)
- Visual Studio Code with Go extension installed globally
- PostgreSQL v13+

### Health Check Endpoint
<PROTOCOL>://<DOMAIN-or-IP>:<PORT>/
### Swagger API Documentation
<PROTOCOL>://<DOMAIN-or-IP>:<PORT>/swagger/index.html

## API Endpoints
- /api/generate [POST]
- /api/resolve [GET]

## Highlights
- High performance, fast APIs with minimal resource foot print
- Data persistence in PostgreSQL
- Open API documentation support using Swagger
- GitHub actions

## TODOs
- HTTPS support / redirection
- PostgreSQL database migration
- Dockerization