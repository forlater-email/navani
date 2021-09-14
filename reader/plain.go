package reader

import (
	"os/exec"
	"strings"
)

func MakePlaintext(html []byte) ([]byte, error) {
	args := []string{"-image_links", "-dump", "-stdin"}
	cmd := exec.Command("lynx", args...)
	cmd.Stdin = strings.NewReader(string(html))
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return out, nil
}
