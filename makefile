postgres:
	docker run --name postgres16 --network=bank-network -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -e POSTGRES_DB=simple_bank -d postgres:16-alpine

stpostgres:
	docker stop postgres16

rmpostgres:
	docker rm postgres16

createdb:
	docker exec -it postgres16 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres16 dropdb simple_bank

migrateup:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5433/simple_bank?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5433/simple_bank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5433/simple_bank?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5433/simple_bank?sslmode=disable" -verbose down 1

migratedrop:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5433/simple_bank?sslmode=disable" -verbose drop

sqlc:
	sqlc generate

test:
	go test -v --cover ./...

server:
	go run main.go

mockgen:
	mockgen -package mockdb -destination db/mock/store.go github.com/thanhphuocnguyen/go-simple-bank/db/sqlc Store

createmigration:
	migrate create -ext sql -dir db/migrations -seq $(name)

.PHONY: postgres createdb dropdb stpostgres rmpostgres migrateup migratedown sqlc migratedrop test server mockgen createmigration migratedown1 migrateup1
