postgres: ; sudo docker run --name postgrescontainr -p 5430:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb: ; sudo docker exec -it postgrescontainr createdb --username=root --owner=root go_bank

dropdb: ; sudo docker exec -it postgrescontainr dropdb go_bank

schema: ; migrate create -ext sql -dir db/migration -seq init_schema

migrateup: ; migrate -path db/migration -database "postgresql://root:secret@localhost:5430/go_bank?sslmode=disable" -verbose up

migratedown: ; migrate -path db/migration -database "postgresql://root:secret@localhost:5430/go_bank?sslmode=disable" -verbose down

migrateup1: ; migrate -path db/migration -database "postgresql://root:secret@localhost:5430/go_bank?sslmode=disable" -verbose up 1

migratedown1: ; migrate -path db/migration -database "postgresql://root:secret@localhost:5430/go_bank?sslmode=disable" -verbose down 1

sqlc: ; sqlc generate

test: ; go test -v -cover ./...

server: ; go run main.go

mock: ;  mockgen -package mockdb -destination db/mock/store.go gobank/db/sqlc Store

.PHONY: ; postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 sqlc test server mock
