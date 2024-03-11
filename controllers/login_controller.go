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

type LoginForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(c *gin.Context, redisClient redis.Conn) {
	repo := models.NewRepository(db.DB)

	var form LoginForm

	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	email := form.Email
	password := form.Password

	errMsg, err := models.LoadConfig("settings/error_messages.json")
	if err != nil {
		fmt.Println("Error loading config:", err)
	}

	u, err := repo.FindUserByEmail(email)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "login/login.html", gin.H{
			"errorMessages": []string{errMsg.LoginError},
			"User":          u,
		})
		return
	}

	if u == nil {
		c.HTML(http.StatusUnauthorized, "login/login.html", gin.H{
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
		c.HTML(http.StatusUnauthorized, "login/login.html", gin.H{
			"errorMessages": []string{errMsg.LoginError},
			"User":          u,
		})
		return
	}

	if !u.IsActive {
		c.HTML(http.StatusUnauthorized, "login/login.html", gin.H{
			"errorMessages": []string{errMsg.InactiveAccount},
			"User":          u,
		})
		return
	}

	c.SetCookie("uid", strconv.FormatUint(u.ID, 10), 3600, "/", "localhost", true, true)

	// ログイン成功
	c.JSON(http.StatusOK, gin.H{
		"message": "ログインに成功しました",
	})
}
