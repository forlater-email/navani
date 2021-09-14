package reader

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"git.icyphox.sh/forlater/navani/cache"
	readability "github.com/go-shiori/go-readability"
)

type Article struct {
	readability.Article
	URL *url.URL
}

func checksum(s []byte) string {
	h := sha1.New()
	h.Write(s)
	b := h.Sum(nil)
	return hex.EncodeToString(b)
}

// Fetches the web page and stores the hash of the URL against
// the response body in cache. Returns an io.Reader.
func Fetch(url string) (io.Reader, error) {
	sum := checksum([]byte(url))
	c, err := cache.NewConn()

	body, err := c.Get(sum)
	// Not in cache.
	if err != nil {
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		buf := bytes.Buffer{}
		// Read into r and write into buf.
		// Cache and return!
		r := io.TeeReader(resp.Body, &buf)
		b, err := io.ReadAll(r)
		if err != nil {
			c.Set(b)
		}
		return &buf, nil
	}

	return strings.NewReader(body), nil
}

// Makes a given html body readable. Returns an error if it
// can't.
func Readable(r io.Reader, u *url.URL) (Article, error) {
	article, err := readability.FromReader(r, u)
	if err != nil {
		return Article{}, fmt.Errorf("failed to parse %s: %v\n", u, err)
	}

	return Article{article, u}, nil
}
