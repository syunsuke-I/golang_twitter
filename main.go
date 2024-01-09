package main

import (
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
	r.GET("/login", controllers.LoginPage)
	r.POST("/login", func(c *gin.Context) {
		controllers.Login(c, redisClient)
	})
	r.GET("/activate", controllers.Activate)
	r.GET("/signup", controllers.SignUp)
	r.POST("/signup", controllers.UserCreate)

	// ログイン後にアクセスされるルートにセッション確認ミドルウェアを適用
	authRequired := r.Group("/")
	//authRequired.Use(utils.SessionAuthMiddleware(redisClient)) // 作業中なので
	{
		authRequired.GET("/home", func(c *gin.Context) {
			controllers.Home(c, redisClient)
		})
		authRequired.POST("/tweet", controllers.TweetCreate)
	}

	r.Run(":8080")
}
