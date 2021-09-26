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

func SendArticle(article *reader.Article, to string, readable bool) error {
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
	if err != nil {
		return fmt.Errorf("making plaintext: %w\n", err)
	}

	email := mail.New()
	email.Encryption = mail.EncryptionTLS
	email.SetFrom(fmt.Sprintf("saved forlater <%s>", EMAIL_FROM))
	email.AddTo(to)
	if readable {
		email.SetSubject(article.Title)
		email.SetBody("text/plain", string(plainContent))
		email.AddAlternative("text/html", string(htmlContent))
	} else {
		email.SetSubject(article.URL.String())
		email.SetBody("text/plain", fmt.Sprintf(
			"We were unable to parse your link: %s",
			article.URL.String(),
		))
	}
	email.Username = EMAIL_USER_SECRET
	email.Password = EMAIL_PASSWORD

	err = email.Send(SMTP_HOST + ":" + SMTP_PORT)
	if err != nil {
		return err
	}
	return nil
}
