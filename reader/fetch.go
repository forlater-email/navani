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
	client := &http.Client{}
	sum := checksum([]byte(url))
	c, err := cache.NewConn()
	if err != nil {
		return nil, fmt.Errorf("cache error: %w\n", err)
	}

	body, err := c.Get(sum)
	// Not in cache.
	if err != nil {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("http error: %w\n", err)
		}

		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.82 Safari/537.36")
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("http client error: %w\n", err)
		}

		buf := bytes.Buffer{}
		// Read into r and write into buf.
		// Cache and return!
		r := io.TeeReader(resp.Body, &buf)
		b, err := io.ReadAll(r)
		if err != nil {
			return nil, fmt.Errorf("io error: %w\n", err)
		}
		_, err = c.Set(sum, b)
		if err != nil {
			return nil, fmt.Errorf("cache error: %w\n", err)
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
		return Article{readability.Article{}, u}, fmt.Errorf("failed to parse %s: %w\n", u, err)
	}

	return Article{article, u}, nil
}
