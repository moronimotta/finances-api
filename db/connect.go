// db/connect.go
package db

import (
	"finances-api/entities"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() (Database, error) {

	db, err := gorm.Open(postgres.Open(os.Getenv("DB_URL")), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	if err := db.AutoMigrate(&entities.Gateway{},
		&entities.Invoices{},
		&entities.Products{},
		&entities.UserProducts{},
		&entities.Transactions{}); err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}

	return &GormDatabase{DB: db}, nil
}
