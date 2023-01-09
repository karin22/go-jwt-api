package orm

import (
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
  )
  
var DB *gorm.DB
var err error

func InitDB() {
	dsn := "host=0.0.0.0 user=postgres password=admin dbname=postgres port=5432 sslmode=disable"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}
	// Auto Migrate
	DB.AutoMigrate(&User{})
}