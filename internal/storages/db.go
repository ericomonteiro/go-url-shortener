package storages

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

// NewDB creates a new PostgreSQL database connection
func NewDB() (*sql.DB, error) {
	connStr := "postgresql://postgres:postgres@localhost:5432/url_shortener?sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Error opening database connection: %v", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Printf("Error pinging database: %v", err)
		return nil, err
	}

	log.Println("Successfully connected to database")
	return db, nil
}
