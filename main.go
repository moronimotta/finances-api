package main

import (
	"finances-api/confs"
	"finances-api/db"
	"finances-api/server"
	"log"
	"time"
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

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered from panic: %v", r)
			}
		}()
		time.Sleep(5 * time.Second)
		log.Println("Starting RabbitMQ server...")
		rabbitServer := server.NewRabbitMQServer(db)
		rabbitServer.Start()
	}()

	// Initialize the server
	server := server.NewServer(db)
	server.Start()

}
