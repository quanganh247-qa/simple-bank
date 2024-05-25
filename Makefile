postgres:
	docker run --name postgres12 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:	
	docker exec -it postgres12 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgres://root:fHWFyt98gPR51h3NxjcroWoIscjt7QOb@dpg-cp649mmn7f5s73a6r8ag-a.oregon-postgres.render.com/simple_bank_7qc2"  -verbose up

migratedown:
	migrate -path db/migration -database "postgres://root:fHWFyt98gPR51h3NxjcroWoIscjt7QOb@dpg-cp649mmn7f5s73a6r8ag-a.oregon-postgres.render.com/simple_bank_7qc2"  -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./... 

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go tutorial.sqlc.dev/app/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test 