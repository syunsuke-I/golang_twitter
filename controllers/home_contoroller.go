package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Home(c *gin.Context) {

	c.HTML(
		http.StatusOK,
		"home/home.html",
		gin.H{},
	)
}
