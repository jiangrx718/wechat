package agent

import (
	"testing"
)

func TestSplitText(t *testing.T) {
	// 中文字符串按 rune 切分
	text := "你好世界今天天气真不错我们去玩吧"
	chunks := splitText(text, 5, 2)
	// size=5 overlap=2 => step=3，预期多块且覆盖完整
	if len(chunks) == 0 {
		t.Fatal("no chunks")
	}
	// 第一块应为前 5 个字符
	if chunks[0] != "你好世界今" {
		t.Fatalf("first chunk = %q", chunks[0])
	}
	// 所有内容应被覆盖（拼接 step 部分含全部字符）
	last := chunks[len(chunks)-1]
	if !contains(last, "玩吧") {
		t.Fatalf("last chunk should contain ending, got %q", last)
	}
}

func TestSplitTextShort(t *testing.T) {
	// 文本短于 size 时应整段返回
	chunks := splitText("短文本", 500, 50)
	if len(chunks) != 1 || chunks[0] != "短文本" {
		t.Fatalf("want single chunk, got %v", chunks)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (indexOf(s, sub) >= 0)
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
