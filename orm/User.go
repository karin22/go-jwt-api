package orm

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"user_name"`
	Password string `json:"-"`
	Fullname string `json:"full_name"`
}
