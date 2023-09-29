package orm

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	DB.AutoMigrate(&Files{})

}
