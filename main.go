package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/syunsuke-I/golang_twitter/controllers"
	database "github.com/syunsuke-I/golang_twitter/db"
	"github.com/syunsuke-I/golang_twitter/utils"
)

func main() {

	r := gin.Default()
	r.Static("/frontend", "./frontend")
	r.LoadHTMLGlob("frontend/templates/**/**")

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

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
