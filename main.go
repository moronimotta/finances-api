package main

import (
	"finances-api/confs"
	"finances-api/db"
	"finances-api/server"
	"log"
)

func main() {

	err := confs.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}
	// Initialize the database connection
	db, err := db.Connect()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Initialize the server
	server := server.NewServer(db)
	server.Start()

}
