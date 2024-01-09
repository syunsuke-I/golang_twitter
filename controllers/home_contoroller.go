package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/syunsuke-I/golang_twitter/models"
)

func Home(c *gin.Context, redisClient redis.Conn) {
	repo := models.NewRepository(db.DB)
	errMsg, err := models.LoadConfig("settings/error_messages.json")
	if err != nil {
		fmt.Println("Error loading config:", err)
	}

	uid, err := c.Cookie("uid")
	if err != nil {
		// セッションIDが見つからない場合はエラーハンドリング
		c.HTML(http.StatusBadRequest, "login/login.html", gin.H{
			"errorMessages": []string{errMsg.LoginRequired},
		})
		return
	}

	userId, err := strconv.ParseUint(uid, 10, 64)
	if err != nil {
		fmt.Println("Error parsing uid:", err)
	}
	tweet := repo.TweetsFind(userId)
	c.HTML(
		http.StatusOK,
		"home/home.html",
		gin.H{
			"tweet": tweet,
		},
	)
}
