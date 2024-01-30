postgres:
	docker run --name simple_bank -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=admin123 -d postgres

createdb: 
	docker exec -it simple_bank createdb --username=postgres --owner=postgres simple_bank

dropdb: 
	docker exec -it simple_bank dropdb -U postgres simple_bank

migrate_up:
	migrate -path db/migration -database "postgresql://postgres:admin123@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrate_down:
	migrate -path db/migration -database "postgresql://postgres:admin123@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/techschool/simplebank/db/sqlc Store

.PHONY: postgres createdb dropdb migrate_up migrate_down sqlc test server mock