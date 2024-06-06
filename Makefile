postgres:
	docker run --name postgres --network bank-network -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -p 5433:5432 -d postgres

createdb:
	docker exec -it postgres createdb --username=root --owner=root simple_bank

dropdb:	
	docker exec -it postgres dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgres://root:secret@localhost:5433/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgres://root:secret@localhost:5433/simple_bank?sslmode=disable" -verbose down

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

proto:
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative --experimental_allow_proto3_optional\
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
	proto/*.proto
	statik -src=./doc/swagger -dest=./doc

evans:
	evans --host localhost --port 9090 -r repl

.PHONY: postgres createdb dropdb migrateup migratedown  migrateup1 migratedown1 sqlc test server mock proto evans