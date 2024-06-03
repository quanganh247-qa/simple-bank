# Build stage
FROM golang:1.22-alpine3.19 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .

EXPOSE 8080 9090
CMD [ "/app/main" ]



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