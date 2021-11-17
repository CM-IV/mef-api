package main

import (
	"database/sql"
	"log"

	"gitea.civdev.rocks/Occidental-Tech/mef-api/api"
	db "gitea.civdev.rocks/Occidental-Tech/mef-api/db/sqlc"
	"gitea.civdev.rocks/Occidental-Tech/mef-api/db/util"
	_ "github.com/lib/pq"
)



func main() {

	config, err := util.LoadConfig(".")
	if err != nil {

		log.Fatal("cannot load config", err)

	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {

		log.Fatal("cannot connect to db", err)

	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)

	if err != nil {

		log.Fatal("cannot start server", err)

	}

}
