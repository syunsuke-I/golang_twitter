package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/syunsuke-I/golang_twitter/models"
	"github.com/syunsuke-I/golang_twitter/utils"
)

func ShowLoginPage(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"login/login.html",
		gin.H{},
	)
}

func ProcessLogin(c *gin.Context, redisClient redis.Conn) {
	repo := models.NewRepository(db.DB)
	email := c.PostForm("email")
	password := c.PostForm("password")

	errMsg, err := models.LoadConfig("settings/error_messages.json")
	if err != nil {
		fmt.Println("Error loading config:", err)
	}

	u, err := repo.FindUserByEmail(email)
	if err != nil {
		c.HTML(http.StatusBadRequest, "login/login.html", gin.H{
			"errorMessages": []string{errMsg.ServerError},
			"User":          u,
		})
		return
	}

	if u == nil {
		c.HTML(http.StatusBadRequest, "login/login.html", gin.H{
			"errorMessages": []string{errMsg.LoginError},
		})
		return
	}

	if err := utils.SetSessionWithUserID(u.ID, redisClient); err != nil {
		// エラー処理
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"Error": "Internal server error"})
		return
	}

	err = models.CompareHashAndPassword(u.Password, password)
	if err != nil {
		c.HTML(http.StatusBadRequest, "login/login.html", gin.H{
			"errorMessages": []string{errMsg.LoginError},
			"User":          u,
		})
		return
	}

	if !u.IsActive {
		c.HTML(http.StatusBadRequest, "login/login.html", gin.H{
			"errorMessages": []string{errMsg.InactiveAccount},
			"User":          u,
		})
		return
	}

	c.SetCookie("uid", strconv.FormatUint(u.ID, 10), 3600, "/", "localhost", true, true)

	// ログイン成功
	c.Redirect(http.StatusMovedPermanently, "home")
}
