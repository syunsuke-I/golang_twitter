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
		"sign_up/sign_up.html",
		gin.H{},
	)
}

func Activate(c *gin.Context) {
	repo := models.NewRepository(db.DB)
	token := c.Query("token")
	u, err := repo.FindUserByActivationToken(token)
	if err != nil {
		log.Fatal(err)
	}

	if u == nil {
		log.Fatal("User not found")
	}

	err = repo.ActivateUser(u)
	if err != nil {
		log.Fatal(err)
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
	sendMail(email, token)
	c.Redirect(http.StatusMovedPermanently, "home")
}

func sendMail(email string, token string) {
	body := createEmailBody(token)
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

	msg.Subject(mime.BEncoding.Encode("UTF-8", "アクティベーションを完了してください。"))
	msg.SetBodyString(mail.TypeTextPlain, body)

	c, err := mail.NewClient(host, mail.WithPort(port))
	c.SetTLSPolicy(mail.TLSOpportunistic)
	if err != nil {
		log.Fatal(err)
	}
	if err := c.DialAndSend(msg); err != nil {
		log.Fatal(err)
	}
}

func createEmailBody(token string) string {
	bodyMsg := `
	こんにちは、

	ご登録ありがとうございます。アカウントのアクティベーションを完了するために、以下のリンクをクリックしてください。

	http://localhost:8080/activate?token=` + token + `

	このリンクは24時間有効です。この期間内にアクティベーションを完了してください。

	もしこのメールに心当たりがない場合は、無視していただいて構いません。

	よろしくお願いします。
	`
	return bodyMsg
}
