
psa:
	docker ps -a

postgres:
	docker run --name postgres15 -p 12346:5432 -e POSTGRES_USER=root  -e POSTGRES_PASSWORD=secret -d postgres:15.2-alpine

createdb:
	docker exec -it postgres15  createdb  --username=root  --owner=root simple_bank

dropdb:
	docker exec -it postgres15  dropdb  simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:12346/simple_bank?sslmode=disable" -verbose up
migrateup1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:12346/simple_bank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:12346/simple_bank?sslmode=disable" -verbose down
migratedown1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:12346/simple_bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...
server:
	go run main.go
.PHONY: createdb dropdb postgres migrateup migratedown sqlc test server migratedown1 migrateup1
