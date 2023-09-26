package auth

import (
	"helloGo/jwt-api/orm"
	"net/http"

	"helloGo/jwt-api/service"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterBody struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Fullname string `json:"fullname" binding:"required"`
}

func Register(c *gin.Context) {
	var json RegisterBody
	// Call BindJSON to bind the received JSON to
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}

	//Check User Already
	var userAlready orm.User
	orm.DB.Where("username = ?", json.Username).Find(&userAlready)

	if userAlready.ID > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "User Already",
		})
		return
	}

	//Create User
	encryptPasswords, _ := bcrypt.GenerateFromPassword([]byte(json.Password), 10)

	user := orm.User{Username: json.Username, Password: string(encryptPasswords), Fullname: json.Fullname}

	orm.DB.Create(&user) // pass pointer of data to Create

	if user.ID > 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "User created Success",
			"userID":  user.ID,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "User created Failed",
		})
	}
}

type LoginBody struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var json LoginBody
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
			"status":  http.StatusBadRequest,
			"message": "User Dose Not Exists",
		})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(userExist.Password), []byte(json.Password))
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "Login Success",
			"token":   service.GenerateToken(userExist.ID),
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Login Failed",
		})
	}
}
