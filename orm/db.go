package orm

import (
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
	"os"
  )
  
var DB *gorm.DB
var err error

func InitDB() {
	dsn := os.Getenv("POSTGRES_DNS")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}
	// Auto Migrate
	DB.AutoMigrate(&User{})
}