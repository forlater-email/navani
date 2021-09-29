package mail

import (
	"fmt"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"mvdan.cc/xurls/v2"
)

type Mail struct {
	From    string
	Date    string
	ReplyTo string
	Body    string
	Parts   map[string]string
}

// TODO
// Strips the signature from an email.
func stripSignature(text string) string {
	lines := strings.Split(text, "\n")
	body := []string{}
	for i := range lines {
		body = append(body, lines[i])
		if lines[i] == "--" || lines[i] == "-- " {
			break
		}
	}
	return strings.Join(body, "\n")
}

// Extracts URLs from given text.
func ExtractURLs(text string) []string {
	x := xurls.Strict()
	return x.FindAllString(text, -1)
}

// Returns the main body of the email; hopefully containing URLs.
func MailBody(parts map[string]string) (string, error) {
	if plain, ok := parts["text/plain"]; ok {
		return stripSignature(plain), nil
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
