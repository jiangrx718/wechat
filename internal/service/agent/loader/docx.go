package loader

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

// loadDocx 从 docx 字节流中提取纯文本。
//
// docx 本质是一个 zip 包，正文位于 word/document.xml。
// 文本被包裹在 <w:t> 标签中，段落由 <w:p> 分隔。
func loadDocx(r io.Reader) (string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return "", fmt.Errorf("读取 docx 失败: %w", err)
	}

	zipReader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", fmt.Errorf("打开 docx zip 失败: %w", err)
	}

	var docXML []byte
	for _, f := range zipReader.File {
		if f.Name == "word/document.xml" {
			rc, err := f.Open()
			if err != nil {
				return "", fmt.Errorf("打开 document.xml 失败: %w", err)
			}
			docXML, err = io.ReadAll(rc)
			rc.Close()
			if err != nil {
				return "", fmt.Errorf("读取 document.xml 失败: %w", err)
			}
			break
		}
	}
	if docXML == nil {
		return "", fmt.Errorf("docx 内未找到 word/document.xml")
	}

	return parseDocumentXML(docXML), nil
}

// parseDocumentXML 解析 document.xml，按段落拼接文本。
//
// document.xml 结构形如：
//   <w:p>            段落
//     <w:r>          文本块
//       <w:t>文本</w:t>
//     </w:r>
//   </w:p>
func parseDocumentXML(data []byte) string {
	dec := xml.NewDecoder(bytes.NewReader(data))
	var (
		sb       strings.Builder
		inPara   bool // 是否在 <w:p> 段落中
		inText   bool // 是否在 <w:t> 文本节点中
		hasText  bool // 当前段落是否已产出文本（用于决定是否换行）
	)

	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			break // 容错：遇到畸形 XML 直接返回已提取内容
		}

		switch el := tok.(type) {
		case xml.StartElement:
			switch el.Name.Local {
			case "p":
				inPara = true
				hasText = false
			case "t":
				inText = true
			}
		case xml.CharData:
			if inText && inPara {
				sb.Write(el)
				hasText = true
			}
		case xml.EndElement:
			switch el.Name.Local {
			case "t":
				inText = false
			case "p":
				// 段落结束，追加换行
				if hasText {
					sb.WriteByte('\n')
				}
				inPara = false
			}
		}
	}

	return strings.TrimSpace(sb.String())
}
