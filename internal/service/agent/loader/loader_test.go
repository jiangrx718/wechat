package loader

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"
)

func TestLoadText(t *testing.T) {
	got, err := Load("a.txt", strings.NewReader("  hello world  "))
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got != "hello world" {
		t.Fatalf("want %q, got %q", "hello world", got)
	}

	// markdown
	got, err = Load("a.md", strings.NewReader("# 标题\n正文"))
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got != "# 标题\n正文" {
		t.Fatalf("got %q", got)
	}
}

func TestLoadUnsupported(t *testing.T) {
	if _, err := Load("a.xlsx", strings.NewReader("x")); err == nil {
		t.Fatal("want error for unsupported type")
	}
}

func TestParseDocumentXML(t *testing.T) {
	xml := `<w:document xmlns:w="x"><w:body>
		<w:p><w:r><w:t>Hello</w:t></w:r></w:p>
		<w:p><w:r><w:t>World</w:t></w:r><w:r><w:t>!</w:t></w:r></w:p>
	</w:body></w:document>`
	got := parseDocumentXML([]byte(xml))
	// 两段，第二段两个 run 合并
	want := "Hello\nWorld!"
	if got != want {
		t.Fatalf("want %q, got %q", want, got)
	}
}

func TestLoadDocx(t *testing.T) {
	// 构造一个最小 docx zip 包
	xml := `<?xml version="1.0"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
	<w:body>
		<w:p><w:r><w:t>第一个段落</w:t></w:r></w:p>
		<w:p><w:r><w:t>second para</w:t></w:r></w:p>
	</w:body>
</w:document>`

	buf := &bytes.Buffer{}
	w := zip.NewWriter(buf)
	fw, _ := w.Create("word/document.xml")
	fw.Write([]byte(xml))
	w.Close()

	got, err := Load("test.docx", buf)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if !strings.Contains(got, "第一个段落") || !strings.Contains(got, "second para") {
		t.Fatalf("unexpected docx text: %q", got)
	}
}

func TestExtractTextFromContent(t *testing.T) {
	// 模拟一段未压缩的 PDF content stream
	cs := []byte("BT /F1 12 Tf (Hello) Tj 100 0 Td (World) Tj ET")
	got := extractTextFromContent(cs)
	if !strings.Contains(got, "Hello") || !strings.Contains(got, "World") {
		t.Fatalf("want Hello and World, got %q", got)
	}
}
