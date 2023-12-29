package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/syunsuke-I/golang_twitter/models"
)

func ShowLoginPage(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"login/login.html",
		gin.H{},
	)
}

func ProcessLogin(c *gin.Context) {
	repo := models.NewRepository(db.DB)
	email := c.PostForm("email")
	password := c.PostForm("passsword")

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
	err = models.CompareHashAndPassword(u.Password, password)
	fmt.Println("err = ", err)
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

	// ログイン成功
	c.Redirect(http.StatusMovedPermanently, "home")
}
