package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/syunsuke-I/golang_twitter/controllers"
	database "github.com/syunsuke-I/golang_twitter/db"
	"github.com/syunsuke-I/golang_twitter/utils"
)

func main() {

	r := gin.Default()
	r.Static("/frontend", "./frontend")
	r.LoadHTMLGlob("frontend/templates/**/**")

	// データベース接続の初期化
	db := database.NewDatabase()
	defer db.Close()

	redisClient := utils.RedisConnection()
	defer redisClient.Close()

	if err := db.CreateTables(); err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	controllers.Init(db)

	// ルーティング設定
	r.GET("/", controllers.Top)
	r.GET("/login", controllers.ShowLoginPage)
	r.POST("/login", func(c *gin.Context) {
		controllers.ProcessLogin(c, redisClient)
	})
	r.GET("/activate", controllers.Activate)
	r.GET("/signup", controllers.SignUp)
	r.POST("/signup", controllers.UserCreate)

	// ログイン後にアクセスされるルートにセッション確認ミドルウェアを適用
	authRequired := r.Group("/")
	authRequired.Use(SessionAuthMiddleware())
	{
		authRequired.GET("/home", controllers.Home)
	}

	r.Run(":8080")
}

func SessionAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("hoge")
	}
}
