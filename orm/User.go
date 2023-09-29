package orm

import (
	"helloGo/jwt-api/model"
)

type User struct {
	model.Model
	Username string `json:"user_name"`
	Password string `json:"-"`
	Fullname string `json:"full_name"`
	FileID   uint   `json:"file_id"`
	File     Files  `gorm:"foreignKey:FileID" json:"file"`
}
