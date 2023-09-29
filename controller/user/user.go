package user

import (
	"helloGo/jwt-api/model"
	"helloGo/jwt-api/orm"
	"helloGo/jwt-api/service"

	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
)

type UpdateUserBody struct {
	Username string `json:"username" binding:"required"`
	Fullname string `json:"fullname"  binding:"required"`
	FileID   int    `json:"file_id" binding:"required"`
}

type File struct {
	ID int `json:"id"`
}

//		ReadAllUsers     godoc
//		@Summary		Read All Users
//		@Description	*Authorization
//		@Tags			User
//		@Produce		json
//		@Router			/users [get]
//	 @security ApiKeyAuth
//	 @Success 200 {object} model.Response "OK"
//	 @Failure 400 {object} model.Response "Bad Request"
//	 @Failure 500 {object} model.Response "Internal Server Error"
func ReadAllUsers(c *gin.Context) {

	users := []*orm.User{}

	orm.DB.Preload("File").Find(&users)
	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "Users Read Success",
		Data:    users,
	})

}

// GetProfile     godoc
// @Summary		Get Profile Users
// @Description	*Authorization
// @Tags			User
// @Produce		json
// @Router			/profile [get]
// @security ApiKeyAuth
// @Success 200 {object} model.Response "OK"
// @Failure 400 {object} model.Response "Bad Request"
// @Failure 500 {object} model.Response "Internal Server Error"
func Profile(c *gin.Context) {
	header := c.Request.Header.Get("Authorization")
	tokenString := strings.TrimPrefix(header, "Bearer ")

	claims, err_valid := service.ValidateToken(tokenString)
	if err_valid != nil {
		return
	}

	user := orm.User{}
	id := claims["userID"]

	err := orm.DB.Where("id = ?", id).Preload("File").Find(&user).Error
	if err != nil {
		logrus.Errorf("find error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    user,
	})
}

//			UpdateUser     godoc
//			@Summary		Update Users
//			@Description	*Authorization
//			@Tags			User
//			@Produce		json
//			@Router			/user/{id} [put]
//	 @security ApiKeyAuth
//		 @param id path int true "id"
//		 @param Body body UpdateUserBody false "body"
//		 @Success 200 {object} model.Response "OK"
//		 @Failure 400 {object} model.Response "Bad Request"
//		 @Failure 500 {object} model.Response "Internal Server Error"
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
		"id":      id,
	})
}

// DeleteUser     godoc
// @Summary		delete Users
// @Description	*Authorization
// @Tags		User
// @Produce		json
// @Router	    /user/{id} [delete]
// @security ApiKeyAuth
// @param id path int true "id"
// @Success 200 {object} model.Response "OK"
// @Failure 400 {object} model.Response "Bad Request"
// @Failure 500 {object} model.Response "Internal Server Error"
func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	user := orm.User{}
	err := orm.DB.Where("id = ?", id).Find(&user).Error
	if err != nil {
		logrus.Errorf("find error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	err = orm.DB.Where("id = ?", id).Delete(&user).Error
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
	})
}
