// auth_controller.go
package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	database "github.com/syunsuke-I/golang_twitter/db"
	"github.com/syunsuke-I/golang_twitter/models"
)

type ErrorMsg struct {
	EmailRequired       string `json:"emailRequired"`
	EmailFormat         string `json:"emailFormat"`
	PasswordRequired    string `json:"passwordRequired"`
	PasswordLength      string `json:"passwordLength"`
	PasswordAlphabet    string `json:"passwordAlphabet"`
	PasswordMixedCase   string `json:"passwordMixedCase"`
	PasswordNumber      string `json:"passwordNumber"`
	PasswordSpecialChar string `json:"passwordSpecialChar"`
	EmailInUse          string `json:"emailInUse"`
}

var db *database.Database

func Init(database *database.Database) {
	db = database
}

func SignUp(c *gin.Context) {

	c.HTML(
		http.StatusOK,
		"sign_up/index.html",
		gin.H{},
	)
}

func UserCreate(c *gin.Context) {
	repo := models.NewRepository(db.DB)
	user := models.User{
		Email:    c.PostForm("email"),
		Password: c.PostForm("password"),
	}

	if _, errorMessages := repo.CreateUser(&user); errorMessages != nil {
		messages := []string{errorMessages.Error()}
		c.HTML(http.StatusBadRequest, "sign_up/index.html", gin.H{
			"errorMessages": messages,
			"User":          user,
		})
		return
	}

	c.Redirect(http.StatusMovedPermanently, "signup")
}
