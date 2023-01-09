package auth

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"helloGo/jwt-api/orm"
  )

type RegisterBody struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Fullname string `json:"fullname" binding:"required"`
	Avatar string `json:"avatar" binding:"required"`
}

func Register(c *gin.Context){
	var  json RegisterBody
	  // Call BindJSON to bind the received JSON to
    if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		  })
        return
    }

	//Check User Exists
	var userExist orm.User
	orm.DB.Where("username = ?", json.Username).Find(&userExist)

	if userExist.ID > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"message": "User Exists",
		  })
		  return
	}

	//Create User 
	encryptPasswords,_ := bcrypt.GenerateFromPassword([]byte(json.Password),10)

	user := orm.User{Username: json.Username , Password: string(encryptPasswords), Fullname: json.Fullname, Avatar: json.Avatar}

	orm.DB.Create(&user) // pass pointer of data to Create

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
}


type LoginBody struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context){
	var  json LoginBody
	  // Call BindJSON to bind the received JSON to
    if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		  })
        return
    }

		//Check User Exists
		var userExist orm.User
		orm.DB.Where("username = ?", json.Username).Find(&userExist)
	
		if userExist.ID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": http.StatusBadRequest,
				"message": "User Dose Not Exists",
			  })
			  return
		}

		err := bcrypt.CompareHashAndPassword([]byte(userExist.Password),[]byte(json.Password))
		if err == nil {
			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusOK,
				"message": "Login Success",
			  })
		}else{
			c.JSON(http.StatusBadRequest, gin.H{
				"status": http.StatusBadRequest,
				"message": "Login Failed",
			  })
		}
}