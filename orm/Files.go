package orm

import (
	"helloGo/jwt-api/model"
)

type Files struct {
	model.Model
	Filename string `json:"file_name"`
	Bucket   string `json:"-"`
	Type     string `json:"type"`
	Path     string `json:"-"`
	Url      string `json:"url"`
}
