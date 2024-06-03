postgres:
	docker run --name postgres -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -p 5432:5432 -d postgres

createdb:
	docker exec -it postgres createdb --username=root --owner=root simple_bank

dropdb:	
	docker exec -it postgres dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

migrateup1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5433/simple_bank?sslmode=disable"  -verbose up 1

migratedown1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:543/simple_bank?sslmode=disable"  -verbose down 1
sqlc:
	sqlc generate

test:
	go test -v -cover ./... 

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go tutorial.sqlc.dev/app/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown  migrateup1 migratedown1 sqlc test server mock