network:
	docker network create mef-network

purge:
	docker-compose -f docker-compose-dev.yml down && docker image rmi mef-api_api

postgres:
	docker run --name postgres --network mef-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=postgres -d postgres:13-alpine

composeupdev:
	docker-compose -f docker-compose-dev.yml up

composestartdev:
	docker-compose -f docker-compose-dev.yml start

composestopdev:
	docker-compose -f docker-compose-dev.yml stop	

composedowndev:
	docker-compose -f docker-compose-dev.yml down


composeup:
	docker-compose up

composestart:
	docker-compose start

composestop:
	docker-compose stop	

composedown:
	docker-compose down

createdb:
	docker exec -it mef-api-postgres-1 createdb --username=root --owner=root meforum

dropdb:
	docker exec -it mef-api-postgres-1 dropdb meforum

migrateup:
	migrate -path db/migration -database "postgresql://root:postgres@localhost:5432/meforum?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:postgres@localhost:5432/meforum?sslmode=disable" -verbose down

migrateup1:
	migrate -path db/migration -database "postgresql://root:postgres@localhost:5432/meforum?sslmode=disable" -verbose up 1

migratedown1:
	migrate -path db/migration -database "postgresql://root:postgres@localhost:5432/meforum?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test-insert:
	go test -count=1 -v ./db/sqlc

test:
	go test -v -cover ./db/sqlc ./api ./token ./util

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/CM-IV/mef-api/db/sqlc Store


.PHONY: composeupup composeupdown composeupstart composeupstop createdb dropdb migrateup migratedown  migrateup1 migratedown1 sqlc test server mock postgres network