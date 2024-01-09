package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/syunsuke-I/golang_twitter/models"
)

func TweetCreate(c *gin.Context) {
	repo := models.NewRepository(db.DB)
	content := c.PostForm("content")

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

	tweet := models.Tweet{
		UserID:  userId,
		Content: content,
	}

	_, errorMessages := repo.CreateTweet(&tweet)

	// エラーがある場合はエラーメッセージを返す
	if errorMessages != nil {
		messages := []string{errorMessages.Error()}
		c.HTML(http.StatusBadRequest, "home/home.html", gin.H{
			"errorMessages": messages,
			"Tweet":         tweet,
		})
		return
	}
	c.Redirect(http.StatusMovedPermanently, "home")
}
