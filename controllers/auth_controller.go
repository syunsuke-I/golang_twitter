package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	database "github.com/syunsuke-I/golang_twitter/db"
	"github.com/syunsuke-I/golang_twitter/models"
	"github.com/syunsuke-I/golang_twitter/utils"
)

var db *database.Database

func Init(database *database.Database) {
	db = database
}

func SignUp(c *gin.Context) {

	c.HTML(
		http.StatusOK,
		"sign_up/sign_up.html",
		gin.H{},
	)
}

func Activate(c *gin.Context) {
	errMsg, err := models.LoadConfig("settings/error_messages.json")
	if err != nil {
		fmt.Println("Error loading config:", err)
	}
	repo := models.NewRepository(db.DB)
	token := c.Query("token")
	u, err := repo.FindUserByActivationToken(token)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "login/login.html", gin.H{
			"errorMessages": []string{errMsg.InvalidActivationToken},
		})
		return
	}

	if u == nil {
		c.HTML(http.StatusUnauthorized, "login/login.html", gin.H{
			"errorMessages": []string{errMsg.InvalidActivationToken},
		})
		return
	}

	err = repo.ActivateUser(u)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "login/login.html", gin.H{
			"errorMessages": []string{errMsg.InvalidActivationToken},
		})
		return
	}

	c.Redirect(http.StatusMovedPermanently, "home")
}

func UserCreate(c *gin.Context) {
	repo := models.NewRepository(db.DB)
	email := c.PostForm("email")
	token := models.GenerateTokenFromEmail(email)
	user := models.User{
		Email:           email,
		Password:        c.PostForm("password"),
		ActivationToken: token,
	}

	if _, errorMessages := repo.CreateUser(&user); errorMessages != nil {
		messages := []string{errorMessages.Error()}
		c.HTML(http.StatusBadRequest, "sign_up/sign_up.html", gin.H{
			"errorMessages": messages,
			"User":          user,
		})
		return
	}
	activationGenerator := utils.ActivationEmailGenerator{}
	utils.SendMail(user.Email, "アクティベーションを完了してください。", activationGenerator, token)
	c.Redirect(http.StatusMovedPermanently, "home")
}
