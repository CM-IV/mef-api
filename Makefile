network:
	docker network create mef-network

postgres:
	docker run --name postgres --network mef-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=postgres -d postgres:13-alpine

composeup:
	docker-compose up

composestart:
	docker-compose start

composestop:
	docker-compose stop	

composedown:
	docker-compose down

createdb:
	docker exec -it postgres createdb --username=root --owner=root meforum

dropdb:
	docker exec -it postgres dropdb meforum

migrateup:
	migrate -path db/migration -database "postgresql://root:postgres@localhost:5432/meforum?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:postgres@localhost:5432/meforum?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test-insert:
	go test -count=1 -v ./db/sqlc

test:
	go test -v -cover ./db/sqlc

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/CM-IV/mef-api/db/sqlc Store


.PHONY: composeupup composeupdown composeupstart composeupstop createdb dropdb migrateup migratedown sqlc test server mock postgres network