package mail

import (
	"bytes"
	"html/template"
	"path/filepath"

	"git.icyphox.sh/forlater/navani/reader"
)

func RenderTemplate(file string, article *reader.Article) ([]byte, error) {
	t, err := template.ParseGlob(filepath.Join("templates", "*.tpl"))
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	if err := t.ExecuteTemplate(buf, file, struct {
		Content template.HTML
		Title   string
		Byline  string
		URL     string
	}{
		template.HTML(article.Content),
		article.Title,
		article.Byline,
		article.URL.String(),
	}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
