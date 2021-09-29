package reader

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"mime"
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

type Response struct {
	Body     io.Reader
	MIMEType string
}

// Fetches the web page and stores the hash of the URL against
// the response body in cache. Returns an io.Reader.
func Fetch(url string) (Response, error) {
	client := &http.Client{}
	sum := checksum([]byte(url))
	c, err := cache.NewConn()
	if err != nil {
		return Response{}, fmt.Errorf("cache error: %w\n", err)
	}

	body, err := c.Get(sum)
	// Not in cache.
	if err != nil {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return Response{}, fmt.Errorf("http error: %w\n", err)
		}

		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.82 Safari/537.36")
		resp, err := client.Do(req)
		if err != nil {
			return Response{}, fmt.Errorf("http client error: %w\n", err)
		}

		mt, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
		if err != nil {
			return Response{}, fmt.Errorf("parse mime: %w\n", err)
		}

		// If page isn't text/html, just return the body; no caching.
		if mt != "text/html" {
			if err != nil {
				return Response{}, fmt.Errorf("reading non-html body: %w\n", err)
			}

			return Response{resp.Body, mt}, nil
		}

		buf := bytes.Buffer{}
		// Read into r and write into buf.
		// Cache and return!
		r := io.TeeReader(resp.Body, &buf)
		b, err := io.ReadAll(r)
		if err != nil {
			return Response{}, fmt.Errorf("io error: %w\n", err)
		}
		_, err = c.Set(sum, b)
		if err != nil {
			return Response{}, fmt.Errorf("cache error: %w\n", err)
		}
		return Response{&buf, mt}, nil
	}

	// We can safely assume it's text/html
	return Response{strings.NewReader(body), "text/html"}, nil
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
