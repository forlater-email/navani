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

func htmlMail(article *reader.Article) (*mail.Email, error) {
	htmlContent, err := RenderTemplate("html.tpl", article)
	if err != nil {
		return &mail.Email{}, err
	}

	htmlAbs, err := rel2abs.Convert(htmlContent, article.URL.String())
	if err != nil {
		return &mail.Email{}, err
	}

	plainContent, err := reader.MakePlaintext(htmlAbs)
	if err != nil {
		return &mail.Email{}, fmt.Errorf("making plaintext: %w\n", err)
	}

	email := mail.NewMSG()
	email.SetBodyData(mail.TextPlain, plainContent)
	email.AddAlternative(mail.TextHTML, string(htmlContent))

	return email, nil
}

func plainMail(article *reader.Article) (*mail.Email, error) {
	email := mail.NewMSG()
	email.SetBodyData(mail.TextPlain, []byte(article.TextContent))
	return email, nil
}

func SendArticle(article *reader.Article, mime string, to string, readable bool) error {
	var (
		EMAIL_FROM = os.Getenv("EMAIL_FROM")
		err        error
		email      *mail.Email
	)

	switch mime {
	case "text/html":
		email, err = htmlMail(article)
		if err != nil {
			return err
		}
	case "html":
		// Exception for weird sites
		email, err = htmlMail(article)
		if err != nil {
			return err
		}
	case "text/plain":
		email, err = plainMail(article)
		if err != nil {
			return err
		}
	default:
		readable = false
	}

	email.SetFrom(fmt.Sprintf("saved forlater <%s>", EMAIL_FROM))
	email.AddTo(to)
	if readable {
		if article.Title != "" {
			email.SetSubject(article.Title)
		} else {
			email.SetSubject(article.URL.String())
		}
	} else {
		email.SetSubject("[forlater.email] Unable to read your link")
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
