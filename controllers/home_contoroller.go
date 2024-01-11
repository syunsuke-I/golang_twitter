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
	tweetDto, err := repo.GetTweets(c, userId)
	if err != nil {
		fmt.Println(err)
	}

	// 通常のリンク用のダミーのスライス
	startPage := max(1, min(tweetDto.Page.Number-1, tweetDto.Page.TotalPages-4))
	endPage := min(tweetDto.Page.TotalPages, max(tweetDto.Page.Number+1, 5))
	var displayedPages []int
	for i := startPage; i <= endPage; i++ {
		displayedPages = append(displayedPages, i)
	}

	// 「前へ」リンク用のページ番号
	prevPage := 0
	if tweetDto.Page.Number > 1 {
		prevPage = tweetDto.Page.Number - 1
	}

	// 「次へ」リンク用のページ番号
	nextPage := 0
	if tweetDto.Page.Number < tweetDto.Page.TotalPages {
		nextPage = tweetDto.Page.Number + 1
	}

	c.HTML(
		http.StatusOK,
		"home/home.html",
		gin.H{
			"Tweets":         tweetDto.Tweets,
			"Pages":          tweetDto.Page,
			"displayedPages": displayedPages,
			"prevPage":       prevPage,
			"nextPage":       nextPage,
		},
	)
}
