package main

import (
	"fmt"
	AuthController "helloGo/jwt-api/controller/auth"
	UserController "helloGo/jwt-api/controller/user"

	"helloGo/jwt-api/middleware"
	"helloGo/jwt-api/orm"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Printf("Error loading .env file")
	}

	orm.InitDB()

	r := gin.Default()

	logrus.SetReportCaller(true)

	protected := r.Group("/", middleware.AuthorizationMiddleware)

	r.POST("/register", AuthController.Register)
	r.POST("/login", AuthController.Login)
	protected.GET("/users", UserController.ReadAllUsers)
	protected.GET("/profile", UserController.Profile)
	protected.PUT("/user/:id", UserController.UpdateUser)

	r.Run("localhost:3000") // listen and serve on 0.0.0.0:3000 (for windows "localhost:3000")
}
