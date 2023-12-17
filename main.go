package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

type Status struct {
	Status string `json:"status"`
}

func main() {
	router := gin.Default()

	status := Status{
		Status: "ok",
	}

	router.GET("/health_check", func(ctx *gin.Context) {
		ctx.JSON(200, status)
	})

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Server Run Failed.: ", err)
	}
}
