package db

import (
	"database/sql"
)

//Store will allow DB execute queries and transactions for all functions
//Composition extending struct functionality
type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {

	return &Store{

		db:      db,
		Queries: New(db),
	}

}
