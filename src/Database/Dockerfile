FROM golang:1.16-alpine
WORKDIR /migrationApp
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN export POSTGRES_HOST=postgres-db
RUN export POSTGRES_PORT=5432
RUN export POSTGRES_USER=postgres
RUN export POSTGRES_PASSWORD=postgres
RUN export POSTGRES_DATABASE=gourlshortenerdb		
RUN echo "Executing PostgreSQL database migrations..."
RUN cd /migrationApp/src/Database/Migrations && go run main.go migrate