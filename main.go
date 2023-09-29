package main

import (
	"fmt"
	AuthController "helloGo/jwt-api/controller/auth"
	FilesController "helloGo/jwt-api/controller/files"
	UserController "helloGo/jwt-api/controller/user"

	"helloGo/jwt-api/middleware"
	"helloGo/jwt-api/orm"

	_ "helloGo/jwt-api/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//	@title			Demo login  API
//	@version		1.0

//	@schemes	https http
//	@host		localhost:3000
//	@BasePath	/

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Printf("Error loading .env file")
	}

	orm.InitDB()

	r := gin.Default()
	logrus.SetReportCaller(true)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	protected := r.Group("/", middleware.AuthorizationMiddleware)

	r.POST("/register", AuthController.Register)
	r.POST("/login", AuthController.Login)

	protected.GET("/users", UserController.ReadAllUsers)

	protected.GET("/profile", UserController.Profile)
	protected.PUT("/user/:id", UserController.UpdateUser)
	protected.DELETE("/user/:id", UserController.DeleteUser)

	r.POST("/upload", FilesController.UploadFile)
	r.Run("localhost:3000") // listen and serve on 0.0.0.0:3000 (for windows "localhost:3000")
}
