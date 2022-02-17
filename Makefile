composeup: ; sudo docker-compose up

postgres: ; sudo docker run --name postgrescontainer  --network bank-network -p 5431:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

gobank: ; sudo docker run --name gobankcontainer --network bank-network -p 8081:8081 -e GIN_MODE=release -e DB_SOURCE=postgresql://root:secret@postgrescontainer:5432/go_bank?sslmode=disable gobank:latest

createdb: ; sudo docker exec -it postgrescontainer createdb --username=root --owner=root go_bank

dropdb: ; sudo docker exec -it postgrescontainer dropdb go_bank

schema: ; migrate create -ext sql -dir db/migration -seq init_schema

migrateup: ; migrate -path db/migration -database "postgresql://root:Gt7QMzKuq7Xo0V9WTVZF@go-bank.cqkrvdhslhri.us-east-1.rds.amazonaws.com:5432/go_bank" -verbose up

migratedown: ; migrate -path db/migration -database "postgresql://root:secret@localhost:5431/go_bank?sslmode=disable" -verbose down

migrateup1: ; migrate -path db/migration -database "postgresql://root:secret@localhost:5431/go_bank?sslmode=disable" -verbose up 1

migratedown1: ; migrate -path db/migration -database "postgresql://root:secret@localhost:5431/go_bank?sslmode=disable" -verbose down 1

sqlc: ; sqlc generate

test: ; go test -v -cover ./...

server: ; go run main.go

mock: ;  mockgen -package mockdb -destination db/mock/store.go gobank/db/sqlc Store

.PHONY: ; postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 sqlc test server mock
