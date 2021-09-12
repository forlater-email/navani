package reader

import (
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
		if b, err := io.ReadAll(resp.Body); err == nil {
			c.Set(sum, b)
		}
		return resp.Body, nil
	}
	return strings.NewReader(body), nil
}

// Makes a given html body readable. Returns an error if it
// can't.
func Readable(r io.Reader, u *url.URL) (readability.Article, error) {
	article, err := readability.FromReader(r, u)
	if err != nil {
		return readability.Article{}, fmt.Errorf("failed to parse %s: %v\n", u, err)
	}

	return article, nil
}
