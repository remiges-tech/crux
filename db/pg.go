package pg

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// Connect initializes a connection to the PostgreSQL database.
func Connect(driverName string, dataSourceName string) (*sql.DB, error) {
	// Open a connection to the database
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Println("Successfully connected to the database")

	return db, nil
}
