package controllers

import (
	"log"
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

	if err := models.CreateUser(c.PostForm("email"), c.PostForm("password")); err != nil {
		log.Println(err)
	}

	c.Redirect(
		http.StatusMovedPermanently,
		"signup",
	)
}
