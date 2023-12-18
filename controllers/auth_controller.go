package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SignUp(c *gin.Context) {

	c.HTML(
		http.StatusOK,
		"sign_up/index.html",
		gin.H{},
	)
}

func UserCreate(c *gin.Context) {

	email := c.PostForm("email")
	password := c.PostForm("password")

	fmt.Println(email)
	fmt.Println(password)
}
