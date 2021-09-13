package mail

import (
	"fmt"
	"log"
	"os"

	"github.com/go-shiori/go-readability"
	"github.com/joegrasse/mail"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
}

func SendArticle(article *readability.Article, to string) error {
	var (
		EMAIL_USER_SECRET = os.Getenv("EMAIL_USER_SECRET")
		EMAIL_PASSWORD    = os.Getenv("EMAIL_PASSWORD")
		EMAIL_FROM        = os.Getenv("EMAIL_FROM")
		SMTP_HOST         = os.Getenv("SMTP_HOST")
		SMTP_PORT         = os.Getenv("SMTP_PORT")
	)

	email := mail.New()
	email.Encryption = mail.EncryptionTLS
	email.SetFrom(fmt.Sprintf("saved forlater <%s>", EMAIL_FROM))
	email.AddTo(to)
	email.SetSubject(article.Title)
	email.SetBody("text/plain", article.TextContent)
	email.AddAlternative("text/html", article.Content)
	email.Username = EMAIL_USER_SECRET
	email.Password = EMAIL_PASSWORD

	err := email.Send(SMTP_HOST + ":" + SMTP_PORT)
	if err != nil {
		return err
	}

	return nil
}
