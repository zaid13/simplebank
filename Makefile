
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

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:12346/simple_bank?sslmode=disable" -verbose down
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
.PHONY: createdb dropdb postgres migrateup migratedown sqlc test
