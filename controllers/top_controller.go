package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Top(c *gin.Context) {

	c.HTML(
		http.StatusOK,
		"top/top.html",
		gin.H{},
	)
}
