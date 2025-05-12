package db

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Config struct {
	Host     string `mapstructure:"DB_HOST" json:"db_host" yaml:"db_host"`
	Port     string `mapstructure:"DB_PORT" json:"db_port" yaml:"db_port"`
	Username string `mapstructure:"DB_USERNAME" json:"db_username" yaml:"db_username"`
	Password string `mapstructure:"DB_PASSWORD" json:"db_password" yaml:"db_password"`
	DBName   string `mapstructure:"DB_NAME" json:"db_name" yaml:"db_name"`
	DBURL    string `mapstructure:"DB_URL" json:"db_url" yaml:"db_url"`
}

func (c *Config) LoadConfig() error {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	c.Host = os.Getenv("DB_HOST")
	c.Port = os.Getenv("DB_PORT")
	c.Username = os.Getenv("DB_USERNAME")
	c.Password = os.Getenv("DB_PASSWORD")
	c.DBName = os.Getenv("DB_NAME")
	c.DBURL = os.Getenv("DB_URL")
	return nil
}

type Database interface {
	GetDB() *gorm.DB
}
