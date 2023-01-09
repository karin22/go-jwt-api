package main

import (
  "github.com/gin-gonic/gin"
  "helloGo/jwt-api/orm"
 AuthController "helloGo/jwt-api/controller/auth"
)


func main() {
	orm.InitDB()

  r := gin.Default()
  r.POST("/register", AuthController.Register)
  r.POST("/login", AuthController.Login)
  r.Run("localhost:3000") // listen and serve on 0.0.0.0:3000 (for windows "localhost:3000")
}