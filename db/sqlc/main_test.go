package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/CM-IV/mef-api/db/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries

func TestMain(m *testing.M) {

	config, err := util.LoadConfig("../..")
	if err != nil {

		log.Fatal("cannot load config", err)

	}

	testDB, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {

		log.Fatal("cannot connect to db", err)

	}

	testQueries = New(testDB)

	os.Exit(m.Run())

}
