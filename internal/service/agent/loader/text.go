package loader

import (
	"io"
	"strings"
)

// loadText 读取纯文本文件（txt / markdown），并去除首尾空白。
func loadText(r io.Reader) (string, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(b)), nil
}
