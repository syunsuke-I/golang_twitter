// main.go
package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/syunsuke-I/golang_twitter/controllers"
	database "github.com/syunsuke-I/golang_twitter/db"
)

func main() {

	r := gin.Default()
	r.Static("/frontend", "./frontend")
	r.LoadHTMLGlob("frontend/templates/**/**")

	// データベース接続の初期化
	db := database.NewDatabase()
	defer db.Close()

	if err := db.CreateTables(); err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	controllers.Init(db)

	// ルーティング設定
	r.GET("/", controllers.Top)
	r.GET("/login", controllers.ShowLoginPage)
	r.POST("/login", controllers.ProcessLogin)
	r.GET("/activate", controllers.Activate)
	r.GET("/home", controllers.Home)
	r.GET("/signup", controllers.SignUp)
	r.POST("/signup", controllers.UserCreate)

	r.Run(":8080")
}
