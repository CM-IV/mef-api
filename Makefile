postgres-up:
	docker-compose up

postgres-start:
	docker-compose start

postgres-stop:
	docker-compose stop	

postgres-down:
	docker-compose down

createdb:
	docker exec -it mef-api_db_1 createdb --username=root --owner=root meforum

dropdb:
	docker exec -it mef-api_db_1 dropdb meforum

migrateup:
	migrate -path db/migration -database "postgresql://postgres:PLMqaztgv435@db.fncomjorbfprbxiwqyul.supabase.co:5432/postgres" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://postgres:PLMqaztgv435@db.fncomjorbfprbxiwqyul.supabase.co:5432/postgres" -verbose down

sqlc:
	sqlc generate

test:
	go test -count=1 -v ./db/sqlc

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/CM-IV/mef-api/db/sqlc Store


.PHONY: postgres-up postgres-down postgres-start postgres-stop createdb dropdb migrateup migratedown sqlc test server mock