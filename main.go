package main

import (
	"database/sql"
	"log"

	"gitea.divdev.rocks/Occidental-Tech/mef_api/api"
	db "gitea.divdev.rocks/Occidental-Tech/mef_api/db/sqlc"
	_ "github.com/lib/pq"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:postgres@localhost:5432/meforum?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {

	conn, err := sql.Open(dbDriver, dbSource)

	if err != nil {

		log.Fatal("cannot connect to db", err)

	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)

	if err != nil {

		log.Fatal("cannot start server", err)

	}

}
