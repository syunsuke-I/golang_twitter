package controllers

import (
	"log"
	"mime"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	database "github.com/syunsuke-I/golang_twitter/db"
	"github.com/syunsuke-I/golang_twitter/models"
	"github.com/wneessen/go-mail"
)

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
	email := c.PostForm("email")
	user := models.User{
		Email:    email,
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
	sendMail(email)
	c.Redirect(http.StatusMovedPermanently, "home")
}

func sendMail(email string) {
	msg := mail.NewMsg()
	host := os.Getenv("SMTP_HOST")
	if host == "" {
		log.Fatal("SMTP_HOST required")
	}

	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		log.Fatal(err)
	}
	if err := msg.From("hoge@example.test"); err != nil {
		log.Fatal(err)
	}
	if err := msg.To(email); err != nil {
		log.Fatal(err)
	}

	msg.Subject(mime.BEncoding.Encode("UTF-8", "こんにちはこんにちは"))
	msg.SetBodyString(mail.TypeTextPlain, "ようこそこんにちは")

	c, err := mail.NewClient(host, mail.WithPort(port))
	c.SetTLSPolicy(mail.TLSOpportunistic)
	if err != nil {
		log.Fatal(err)
	}
	if err := c.DialAndSend(msg); err != nil {
		log.Fatal("this ? = ", err)
	}
}
