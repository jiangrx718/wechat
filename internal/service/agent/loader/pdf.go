package loader

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"strings"
)

// loadPDF 从 PDF 字节流中提取纯文本。
//
// 采用纯标准库实现，针对最常见的文本型 PDF：
//  1. 定位所有 stream...endstream 内容块
//  2. 尝试 zlib(FlateDecode) 解压
//  3. 解析文本绘制操作符 Tj / TJ / ' / " 中包含的字符串
//
// 说明：PDF 的文本编码较复杂（字体可能使用自定义编码、CID 等），
// 此实现覆盖大多数标准编码（WinAnsi/ASCII）的文本型 PDF；
// 对扫描件（图片型 PDF）或使用自定义 ToUnicode 的 PDF 可能提取不全。
func loadPDF(r io.Reader) (string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return "", fmt.Errorf("读取 PDF 失败: %w", err)
	}

	streams := extractStreams(data)
	var sb strings.Builder
	for _, raw := range streams {
		decoded, ok := tryInflate(raw)
		if !ok {
			// 未压缩的 content stream，直接使用原始字节
			decoded = raw
		}
		sb.WriteString(extractTextFromContent(decoded))
	}

	return strings.TrimSpace(sb.String()), nil
}

// extractStreams 提取 PDF 中所有 stream...endstream 之间的原始字节。
func extractStreams(data []byte) [][]byte {
	var streams [][]byte
	cursor := 0
	for {
		start := bytes.Index(data[cursor:], []byte("stream"))
		if start < 0 {
			break
		}
		start += cursor + len("stream")
		// stream 关键字后通常跟 CRLF 或 LF
		if start < len(data) && data[start] == '\r' {
			start++
		}
		if start < len(data) && data[start] == '\n' {
			start++
		}

		end := bytes.Index(data[start:], []byte("endstream"))
		if end < 0 {
			break
		}
		streams = append(streams, data[start:start+end])
		cursor = start + end + len("endstream")
	}
	return streams
}

// tryInflate 尝试用 zlib(FlateDecode) 解压数据；失败则返回 (nil, false)。
func tryInflate(data []byte) ([]byte, bool) {
	zr, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, false
	}
	defer zr.Close()
	decoded, err := io.ReadAll(zr)
	if err != nil {
		return nil, false
	}
	return decoded, true
}

// extractTextFromContent 解析单个 content stream，提取文本绘制操作符中的字符串。
//
// 关心的操作符：
//   (s) Tj        显示字符串 s
//   <hex> Tj      显示十六进制字符串
//   [(s1) n (s2)] TJ  显示数组（数字表示字间距调整，可忽略）
//   (s) ' / (s) " 换行后显示字符串
func extractTextFromContent(data []byte) string {
	var (
		sb     strings.Builder
		i      int
		n      = len(data)
		inStr  bool // 在 () 字符串中
		inHex  bool // 在 <> 十六进制串中
		escape bool // 字符串中的转义标记
	)

	// 在文本绘制操作符之后补一个空格，避免相邻词粘连。
	addSeparator := func() {
		if sb.Len() == 0 {
			return
		}
		if last := sb.String()[sb.Len()-1]; last != ' ' && last != '\n' {
			sb.WriteByte(' ')
		}
	}

	for i < n {
		// 字符串内优先处理
		if inStr {
			c := data[i]
			if escape {
				text, advance := decodePDFEscape(data[i:])
				sb.WriteString(text)
				i += advance
				escape = false
				continue
			}
			if c == '\\' {
				escape = true
				i++
				continue
			}
			if c == ')' {
				inStr = false
				i++
				continue
			}
			// 仅保留 ASCII 可打印字符（标准编码），其余忽略以避免乱码
			if c >= 0x20 && c < 0x7f {
				sb.WriteByte(c)
			}
			i++
			continue
		}

		// 十六进制串内处理
		if inHex {
			c := data[i]
			if c == '>' {
				inHex = false
				i++
				continue
			}
			if isHexDigit(c) {
				hi := hexVal(c)
				i++
				for i < n && data[i] == ' ' {
					i++
				}
				lo := byte(0)
				if i < n && isHexDigit(data[i]) {
					lo = hexVal(data[i])
					i++
				}
				b := hi<<4 | lo
				if b >= 0x20 && b < 0x7f {
					sb.WriteByte(b)
				}
				continue
			}
			i++
			continue
		}

		c := data[i]
		switch c {
		case '(':
			inStr = true
			i++
		case '<':
			// << 表示字典开始，跳过；单独 < 为十六进制串
			if i+1 < n && data[i+1] == '<' {
				i += 2
				continue
			}
			inHex = true
			i++
		case '>':
			i++
		default:
			// 识别文本绘制操作符 Tj / TJ / ' / "，命中后补空格分隔
			if matched, advance := matchTextOp(data[i:]); matched {
				addSeparator()
				i += advance
				continue
			}
			i++
		}
	}

	return sb.String()
}

// matchTextOp 检测当前位置是否为文本绘制操作符，返回是否命中及消耗的字节数。
func matchTextOp(data []byte) (bool, int) {
	// ' 和 " 单字符操作符
	if len(data) >= 1 && (data[0] == '\'' || data[0] == '"') {
		return true, 1
	}
	// Tj / TJ 两字符操作符
	if len(data) >= 2 && data[0] == 'T' && (data[1] == 'j' || data[1] == 'J') {
		return true, 2
	}
	return false, 0
}

// decodePDFEscape 解析 PDF 字符串中的转义序列，返回解码文本与消耗的字节数。
// 支持：\n \r \t \b \f \( \) \\ \ddd(八进制)
func decodePDFEscape(data []byte) (string, int) {
	if len(data) == 0 {
		return "", 0
	}
	switch data[0] {
	case 'n':
		return "\n", 1
	case 'r':
		return "\r", 1
	case 't':
		return "\t", 1
	case 'b':
		return "\b", 1
	case 'f':
		return "\f", 1
	case '(', ')', '\\':
		return string(data[0]), 1
	case '\n':
		return "", 1
	case '\r':
		if len(data) > 1 && data[1] == '\n' {
			return "", 2
		}
		return "", 1
	}
	// 八进制转义 \ddd（1~3 位）
	if data[0] >= '0' && data[0] <= '7' {
		val := byte(0)
		count := 0
		for count < 3 && count < len(data) && data[count] >= '0' && data[count] <= '7' {
			val = val*8 + (data[count] - '0')
			count++
		}
		if val >= 0x20 && val < 0x7f {
			return string(val), count
		}
		return "", count
	}
	return "", 1
}

func isHexDigit(c byte) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

func hexVal(c byte) byte {
	switch {
	case c >= '0' && c <= '9':
		return c - '0'
	case c >= 'a' && c <= 'f':
		return c - 'a' + 10
	default:
		return c - 'A' + 10
	}
}
