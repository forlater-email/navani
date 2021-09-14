package reader

import (
	"os/exec"
	"strings"
)

func MakePlaintext(html []byte) ([]byte, error) {
	args := []string{"-T", "text/html", "-dump"}
	cmd := exec.Command("w3m", args...)
	cmd.Stdin = strings.NewReader(string(html))
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return out, nil
}
