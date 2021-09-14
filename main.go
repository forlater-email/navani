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
		if err != nil {
			log.Println(err)
		}

		for _, u := range distinct(mail.ExtractURLs(body)) {
			parsedURL, err := url.Parse(u)
			if err != nil {
				log.Println(err)
			}

			f, err := reader.Fetch(parsedURL.String())
			if err != nil {
				log.Println(err)
			}

			article, err := reader.Readable(f, parsedURL)
			if err != nil {
				log.Println(err)
			}

			err = mail.SendArticle(&article, m.From)
			if err != nil {
				log.Println(err)
			}
		}
		log.Printf("sent mail to %s\n", m.From)
		w.WriteHeader(204)
	})

	http.ListenAndServe(":8001", nil)
}
