package db

import (
	"gorm.io/gorm"
)

type GormDatabase struct {
	DB *gorm.DB
}

func (g *GormDatabase) GetDB() *gorm.DB {
	return g.DB
}
