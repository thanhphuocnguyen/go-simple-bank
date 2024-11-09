postgres:
	docker run --name postgres16 -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16-alpine

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

migratedown:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5433/simple_bank?sslmode=disable" -verbose down

migratedrop:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5433/simple_bank?sslmode=disable" -verbose drop

sqlc:
	sqlc generate

test:
	go test -v --cover ./...

.PHONY: postgres createdb dropdb stpostgres rmpostgres migrateup migratedown sqlc migratedrop test
