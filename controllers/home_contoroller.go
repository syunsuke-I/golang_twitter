package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/syunsuke-I/golang_twitter/models"
	"github.com/syunsuke-I/golang_twitter/utils"
)

func Home(c *gin.Context, redisClient redis.Conn) {

	errMsg, err := models.LoadConfig("settings/error_messages.json")
	if err != nil {
		fmt.Println("Error loading config:", err)
	}

	email, err := c.Cookie("email")
	if err != nil {
		// セッションIDが見つからない場合はエラーハンドリング
		c.HTML(http.StatusBadRequest, "login/login.html", gin.H{
			"errorMessages": []string{errMsg.LoginRequired},
		})
		return
	}

	repo := models.NewRepository(db.DB)
	u, err := repo.FindUserByEmail(email)
	if err != nil {
		c.HTML(http.StatusBadRequest, "login/login.html", gin.H{
			"errorMessages": []string{errMsg.LoginRequired},
		})
		return
	}

	// ここでredisとの照会を行う
	sid, err := utils.GetSessionIDByUserID(u.ID, redisClient)
	if err != nil {
		c.HTML(http.StatusBadRequest, "login/login.html", gin.H{
			"errorMessages": []string{errMsg.SessionInvalid},
			"User":          u,
		})
		return
	}

	err = utils.IsSessionValid(sid, redisClient)
	if err != nil {
		c.HTML(http.StatusBadRequest, "login/login.html", gin.H{
			"errorMessages": []string{errMsg.SessionInvalid},
			"User":          u,
		})
		return
	}

	c.HTML(
		http.StatusOK,
		"home/home.html",
		gin.H{
			"Email": email,
		})
}
