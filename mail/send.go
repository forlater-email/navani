package mail

import (
	"fmt"
	"log"
	"os"

	"git.icyphox.sh/forlater/navani/reader"
	"github.com/joegrasse/mail"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
}

func SendArticle(article *reader.Article, to string) error {
	var (
		EMAIL_USER_SECRET = os.Getenv("EMAIL_USER_SECRET")
		EMAIL_PASSWORD    = os.Getenv("EMAIL_PASSWORD")
		EMAIL_FROM        = os.Getenv("EMAIL_FROM")
		SMTP_HOST         = os.Getenv("SMTP_HOST")
		SMTP_PORT         = os.Getenv("SMTP_PORT")
	)

	htmlContent, err := RenderTemplate("html.tpl", article)
	if err != nil {
		return err
	}

	plainContent, err := reader.MakePlaintext(htmlContent)

	email := mail.New()
	email.Encryption = mail.EncryptionTLS
	email.SetFrom(fmt.Sprintf("saved forlater <%s>", EMAIL_FROM))
	email.AddTo(to)
	email.SetSubject(article.Title)
	email.SetBody("text/plain", string(plainContent))
	email.AddAlternative("text/html", string(htmlContent))
	email.Username = EMAIL_USER_SECRET
	email.Password = EMAIL_PASSWORD

	err = email.Send(SMTP_HOST + ":" + SMTP_PORT)
	if err != nil {
		return err
	}

	return nil
}
