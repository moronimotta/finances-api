// db/connect.go
package db

import (
	"finances-api/entities"
	"log"
	"log/slog"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() (Database, error) {

	db, err := gorm.Open(postgres.Open(os.Getenv("DB_URL")), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
		slog.Error("Error connecting to database", err)
	}

	if err := db.AutoMigrate(&entities.Gateway{},
		&entities.Invoices{},
		&entities.Products{},
		&entities.UserProducts{},
		&entities.Transactions{},
		&entities.TransactionItem{},
	); err != nil {
		log.Fatalf("Error migrating database: %v", err)
		slog.Error("Error migrating database", err)
	}

	return &GormDatabase{DB: db}, nil
}
