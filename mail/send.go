package mail

import (
	"fmt"
	"io"
	"log"
	"mime"
	"os"

	"git.icyphox.sh/forlater/navani/reader"
	"git.icyphox.sh/rel2abs"
	"github.com/joho/godotenv"
	mail "github.com/xhit/go-simple-mail/v2"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
}

func SendArticle(article *reader.Article, to string, readable bool) error {
	var EMAIL_FROM = os.Getenv("EMAIL_FROM")
	htmlContent, err := RenderTemplate("html.tpl", article)
	if err != nil {
		return err
	}

	htmlAbs, err := rel2abs.Convert(htmlContent, article.URL.String())
	if err != nil {
		htmlAbs = htmlContent
	}

	plainContent, err := reader.MakePlaintext(htmlAbs)
	if err != nil {
		return fmt.Errorf("making plaintext: %w\n", err)
	}

	email := mail.NewMSG()
	email.SetFrom(fmt.Sprintf("saved forlater <%s>", EMAIL_FROM))
	email.AddTo(to)
	if readable {
		email.SetSubject(article.Title)
		email.SetBodyData(mail.TextPlain, plainContent)
		email.AddAlternative(mail.TextHTML, string(htmlContent))
	} else {
		email.SetSubject(article.URL.String())
		email.SetBody(mail.TextPlain, fmt.Sprintf(
			"We were unable to parse your link: %s",
			article.URL.String(),
		))
	}

	c, err := mailClient()
	if err != nil {
		return fmt.Errorf("mail: %w\n", err)
	}
	err = email.Send(c)
	if err != nil {
		return err
	}
	return nil
}

func SendAttachment(r reader.Response, to string, url string) error {
	var EMAIL_FROM = os.Getenv("EMAIL_FROM")
	email := mail.NewMSG()
	email.SetFrom(fmt.Sprintf("saved forlater <%s>", EMAIL_FROM))
	email.AddTo(to)

	b, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("read attachment: %w\n", err)
	}
	fmt.Println(len(b))

	ext, _ := mime.ExtensionsByType(r.MIMEType)
	var name string
	if ext != nil {
		name = "file" + ext[0]
	} else {
		name = "file"
	}
	email.SetSubject(url)
	email.Attach(&mail.File{MimeType: r.MIMEType, Data: b, Name: name, Inline: true})
	email.SetBody(mail.TextPlain, fmt.Sprintf(`That didn't look like a HTML page; we found %s.
We've attached it to this email.`, r.MIMEType))

	c, err := mailClient()
	if err != nil {
		return fmt.Errorf("mail: %w\n", err)
	}
	err = email.Send(c)
	if err != nil {
		return err
	}
	return nil
}
