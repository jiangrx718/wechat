package loader

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

// Load 根据文件名后缀分发到对应的解析器，返回文件的纯文本内容。
//
// 支持的格式：
//   - txt / markdown: 直接读取
//   - pdf: 解析 PDF 文本流
//   - docx: 解析 word 文档
func Load(filename string, r io.Reader) (string, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".txt", ".md", ".markdown", "":
		return loadText(r)
	case ".pdf":
		return loadPDF(r)
	case ".docx":
		return loadDocx(r)
	default:
		return "", fmt.Errorf("unsupported file type: %s", ext)
	}
}
