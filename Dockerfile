FROM golang:1.16-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN cd /app/src && go build -o main .
EXPOSE 3356
ENTRYPOINT [ "./src/main" ]