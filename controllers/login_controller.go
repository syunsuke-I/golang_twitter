package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {

	c.HTML(
		http.StatusOK,
		"login/login.html",
		gin.H{},
	)
}
