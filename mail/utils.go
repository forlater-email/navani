package mail

import (
	"fmt"

	"github.com/microcosm-cc/bluemonday"
	"mvdan.cc/xurls/v2"
)

type Mail struct {
	From    string
	Date    string
	ReplyTo string
	Parts   map[string]string
}

// Extracts URLs from given text.
func ExtractURLs(text string) []string {
	x := xurls.Strict()
	return x.FindAllString(text, -1)
}

// Returns the main body of the email; hopefully containing URLs.
func MailBody(parts map[string]string) (string, error) {
	if plain, ok := parts["text/plain"]; ok {
		return plain, nil
	} else if html, ok := parts["text/html"]; ok {
		p := bluemonday.NewPolicy()
		p.AllowStandardURLs()
		p.AddSpaceWhenStrippingTag(true)
		clean := p.Sanitize(html)
		return clean, nil
	} else {
		return "", fmt.Errorf("no good MIME type found")
	}
}
