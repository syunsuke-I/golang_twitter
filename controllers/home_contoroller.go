package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
)

func Home(c *gin.Context, redisClient redis.Conn) {
	c.HTML(
		http.StatusOK,
		"home/home.html",
		gin.H{})
}
