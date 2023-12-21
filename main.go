package main

import (
	"github.com/gin-gonic/gin"

	"github.com/syunsuke-I/golang_twitter/controllers"
	"github.com/syunsuke-I/golang_twitter/models"

	_ "github.com/lib/pq"
)

func main() {

	// ルーティング
	router := gin.Default()
	router.Static("/frontend", "./frontend")
	router.LoadHTMLGlob("frontend/templates/**/**")
	router.GET("/signup", controllers.SignUp)
	router.POST("/signup", controllers.UserCreate)

	// DBの初期化
	models.Init()

	router.Run("0.0.0.0:8080")
}
