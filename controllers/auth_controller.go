package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/syunsuke-I/golang_twitter/models"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init() {
	var err error
	db, err = models.OpenDatabaseConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	repo := models.NewRepository(db)
	repo.CreateTables()
}

func SignUp(c *gin.Context) {

	c.HTML(
		http.StatusOK,
		"sign_up/index.html",
		gin.H{},
	)
}

func UserCreate(c *gin.Context) {

	repo := models.NewRepository(GetDB())

	user := models.User{
		Email:    c.PostForm("email"),
		Password: c.PostForm("password"),
	}

	if _, errorMessages := repo.CreateUser(&user); errorMessages != nil {
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

func GetDB() *gorm.DB {
	return db
}
