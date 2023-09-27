package auth

import (
	"helloGo/jwt-api/orm"
	"net/http"

	"helloGo/jwt-api/service"

	"helloGo/jwt-api/model"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterBody struct {
	Username string `json:"username" binding:"required" example:"karin"`
	Password string `json:"password" binding:"required" example:"1234"`
	Fullname string `json:"fullname" binding:"required" example:"karin pimloy"`
}

// Register     godoc
// @Summary		Register
// @Description
// @Tags		Auth
// @Produce		json
// @Router		/register [post]
// @param Body body RegisterBody false "body"
// @Success 200 {object} model.ResponseLogin "OK"
// @Failure 400 {object} model.Response "Bad Request"
// @Failure 500 {object} model.Response "Internal Server Error"
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
		response := model.ResponseLogin{}
		response.Status = http.StatusOK
		response.Message = "User created Success"
		response.Data.Token = service.GenerateToken(user.ID)

		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "User created Failed",
		})
	}
}

type LoginBody struct {
	Username string `json:"username" binding:"required" example:"karin5"`
	Password string `json:"password" binding:"required" example:"1234"`
}

// Login     godoc
// @Summary		Register
// @Description
// @Tags		Auth
// @Produce		json
// @Router		/login [post]
// @param Body body LoginBody false "body"
// @Success 200 {object} model.ResponseLogin "OK"
// @Failure 400 {object} model.Response "Bad Request"
// @Failure 500 {object} model.Response "Internal Server Error"
func Login(c *gin.Context) {
	json := LoginBody{}
	// Call BindJSON to bind the received JSON to
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}

	//Check User Exists
	userExist := orm.User{}
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

		response := model.ResponseLogin{}
		response.Status = http.StatusOK
		response.Message = "Login Success"
		response.Data.Token = service.GenerateToken(userExist.ID)

		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Login Failed",
		})
	}
}
