package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"git.icyphox.sh/forlater/navani/mail"
	"git.icyphox.sh/forlater/navani/reader"
)

func main() {
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		m := mail.Mail{}
		json.NewDecoder(r.Body).Decode(&m)
		body, err := mail.MailBody(m.Parts)
		log.Printf("recieved webhook: %v\n", m.From)
		if err != nil {
			log.Printf("using body as is: %v\n", err)
			body = mail.StripSignature(m.Body)
		}

		urls := mail.ExtractURLs(body)
		if len(urls) == 0 {
			log.Printf("no urls found")
		}
		for _, u := range distinct(urls) {
			log.Printf("url: %s\n", u)
			parsedURL, err := url.Parse(u)
			if err != nil {
				log.Printf("url parse: %v\n", err)
			}

			resp, err := reader.Fetch(parsedURL.String())
			if err != nil {
				log.Printf("reader fetch: %v\n", err)
			}

			article, err := reader.Readable(resp.Body, parsedURL)
			if (err == nil) && (resp.MIMEType == "text/html") {
				err = mail.SendArticle(&article, m.From, true)
				if err != nil {
					log.Printf("error sending mail to: %s: %v\n", m.From, err)
				} else {
					log.Printf("sent mail to %s: %s\n", m.From, article.Title)
				}
			} else {
				log.Printf("not readable: %s: %s\n", article.URL.String(), resp.MIMEType)
				err := mail.SendArticle(&article, m.From, false)
				if err != nil {
					log.Printf("error sending mail to: %s: %v\n", m.From, err)
				}
			}
		}
		w.WriteHeader(204)
	})

	http.ListenAndServe(":8001", nil)
}
