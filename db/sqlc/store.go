package db

import (
	"database/sql"
)

type Store interface {
	Querier
}

//Store will allow DB execute queries and transactions for all functions
//Composition extending struct functionality
type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {

	return &SQLStore{

		db:      db,
		Queries: New(db),
	}

}
