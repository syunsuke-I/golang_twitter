package main

import (
	"github.com/gin-gonic/gin"

	"github.com/syunsuke-I/golang_twitter/controllers"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("frontend/templates/**/**")
	router.GET("/signup", controllers.SignUp)
	router.Run("0.0.0.0:8080")
}
