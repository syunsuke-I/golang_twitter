package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/syunsuke-I/golang_twitter/models"
)

func SignUp(c *gin.Context) {

	c.HTML(
		http.StatusOK,
		"sign_up/index.html",
		gin.H{},
	)
}

func UserCreate(c *gin.Context) {

	user := models.User{
		Email:    c.PostForm("email"),
		Password: c.PostForm("password"),
	}

	if _, errorMessages := models.CreateUser(&user); errorMessages != nil {
		// エラーメッセージを取得
		messages := []string{errorMessages.Error()}

		// 同じサインアップページを再度レンダリングし、エラーメッセージを渡す
		c.HTML(http.StatusBadRequest, "sign_up/index.html", gin.H{
			"errorMessages": messages,
			"User":          user,
		})
		return
	}

	c.Redirect(
		http.StatusMovedPermanently,
		"signup",
	)
}
