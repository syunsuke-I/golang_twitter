package utils

import (
	"log"
	"mime"
	"os"
	"strconv"

	"github.com/wneessen/go-mail"
)

// メール本文の生成する関数を抽象化するため
type EmailBodyGenerator interface {
	GenerateBody(data interface{}) string
}

type ActivationEmailGenerator struct{}

func (gen ActivationEmailGenerator) GenerateBody(data interface{}) string {
	token, ok := data.(string)
	if !ok {
		return "データの形式が正しくありません。"
	}
	bodyMsg := `
	こんにちは、

	ご登録ありがとうございます。アカウントのアクティベーションを完了するために、以下のリンクをクリックしてください。

	http://localhost:8080/activate?token=` + token + `

	もしこのメールに心当たりがない場合は、無視していただいて構いません。

	よろしくお願いします。
	`

	return bodyMsg
}

// @param emailAddress {string}             送信先のメールアドレス
// @param subject      {string}             メールの件名
// @param generator    {EmailBodyGenerator} メール本文を生成するための EmailBodyGenerator インターフェースを実装したオブジェクト
// @param data         {interface{}}        メール本文の生成に使用されるデータ

func SendMail(emailAddress string, subject string, generator EmailBodyGenerator, data interface{}) {
	body := generator.GenerateBody(data)
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
	if err := msg.To(emailAddress); err != nil {
		log.Fatal(err)
	}

	msg.Subject(mime.BEncoding.Encode("UTF-8", subject))
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
