package user

import (
	"helloGo/jwt-api/orm"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"helloGo/jwt-api/service"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type UpdateUserBody struct {
	Username string `json:"username" binding:"required" `
	Fullname string `json:"fullname"  binding:"required"`
}

type User struct {
	gorm.Model
	Username string
	Fullname string
}

func ReadAllUsers(c *gin.Context) {

	users := []*orm.User{}

	orm.DB.Find(&users)

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Users Read Success",
		"users":   users,
	})

}

func Profile(c *gin.Context) {
	header := c.Request.Header.Get("Authorization")
	tokenString := strings.TrimPrefix(header, "Bearer ")

	claims, err_valid := service.ValidateToken(tokenString)
	if err_valid != nil {
		return
	}

	user := orm.User{}
	id := claims["userID"]
	err := orm.DB.Where("id = ?", id).Find(&user).Error
	if err != nil {
		logrus.Errorf("find error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Success",
		"users":   user,
	})
}

func UpdateUser(c *gin.Context) {
	userBody := UpdateUserBody{}

	if err := c.BindJSON(&userBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}

	user := orm.User{}
	id := c.Param("id")

	err := orm.DB.Where("id = ?", id).Find(&user).Error
	if err != nil {
		logrus.Errorf("find error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	copier.Copy(&user, &userBody)

	err = orm.DB.Where("id = ?", id).Updates(user).Error
	if err != nil {
		logrus.Errorf("find error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Success",
		"users":   user,
	})
}

// func ReadUser(c *gin.Context) {
// 	id := c.Param("id")

// 	var userExist orm.User
// 	orm.DB.Where("id = ?", json.Username).Find(&userExist)
// 	orm.DB.Find(&users)

// 	c.JSON(http.StatusOK, gin.H{
// 		"status":  http.StatusOK,
// 		"message": "Users Read Success",
// 		"users":   users,
// 	})
// }
