package main

import (
	"finances-api/db"
	"finances-api/server"
	"log"
)

func main() {
	// Initialize the database connection
	db, err := db.Connect()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Initialize the server
	server := server.NewServer(db)
	server.Start()

}
