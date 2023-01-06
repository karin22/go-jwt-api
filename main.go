package main

import (
  "net/http"
  "gorm.io/gorm"
  "gorm.io/driver/postgres"
  "github.com/gin-gonic/gin"
  "golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Username  string
	Password string
	Fullname string
	Avatar string
  }

type RegisterBody struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Fullname string `json:"fullname" binding:"required"`
	Avatar string `json:"avatar" binding:"required"`
}
func main() {
	dsn := "host=0.0.0.0 user=postgres password=admin dbname=postgres port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}
	// Auto Migrate
	db.AutoMigrate(&User{})


  r := gin.Default()
  r.POST("/register", func(c *gin.Context) {
	var  json RegisterBody
	  // Call BindJSON to bind the received JSON to
    if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		  })
        return
    }

	//Check User Exists
	var userExist User
	db.Where("username = ?", json.Username).Find(&userExist)

	if userExist.ID > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"message": "User Exists",
		  })
		  return
	}

	//Create User 
	encryptPasswords,_ := bcrypt.GenerateFromPassword([]byte(json.Password),10)

	user := User{Username: json.Username , Password: string(encryptPasswords), Fullname: json.Fullname, Avatar: json.Avatar}

	db.Create(&user) // pass pointer of data to Create

	if user.ID > 0 {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"message": "User created Success",
			"userID" : user.ID,
		  })
	}else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"message": "User created Failed",
			"userID" : user.ID,
		  })
	}

  
  })
  r.Run("localhost:3000") // listen and serve on 0.0.0.0:3000 (for windows "localhost:3000")
}