# Build stage
FROM golang:1.22-alpine3.19 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz

# Run stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./migration

EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT ["/app/start.sh"]



#docker build -t simplebank:latest .

### Run docker image with same port simplebank
# docker run --name simplebank -p 8080:8080 -e GIN_MODE=release -e DB_SOURCE="postgres://root:secret@127.17.0.3:5433/simple_bank?sslmode=disable" simplebank:latest 
### Read IP address of docker images
# docker container inspect postgres
# docker container inspect simplebank

### Create netwok
# docker network create bank-network

### Make postgres use bank-network
# docker network connect bank-network postgres

# docker run --name simplebank --network bank-network -p 8080:8080 -e GIN_MODE=release -e DB_SOURCE="postgres://root:secret@postgres:5432/simple_bank?sslmode=disable" simplebank:latest

## chmod +x ...sh (permission)