include app.env

build:
	docker build -t mef:latest .

run:
	docker run --name mef -p 8080:8080 mef:latest

composeup:
	docker-compose up

composestart:
	docker-compose start

composestop:
	docker-compose stop	

composedown:
	docker-compose down

createdb:
	docker exec -it mef-api_db_1 createdb --username=root --owner=root meforum

dropdb:
	docker exec -it mef-api_db_1 dropdb meforum

migrateup:
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose down

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


.PHONY: composeupup composeupdown composeupstart composeupstop createdb dropdb migrateup migratedown sqlc test server mock