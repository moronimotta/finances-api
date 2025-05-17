// db/connect.go
package db

import (
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

	// if err := db.AutoMigrate(&user.User{}); err != nil {
	// 	log.Fatalf("Error migrating database: %v", err)
	// }

	return &GormDatabase{DB: db}, nil
}
