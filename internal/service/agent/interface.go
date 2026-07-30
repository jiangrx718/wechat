package agent

import (
	"context"
	"wechat-tools/internal/common"
)

// AgentServiceIFace 知识库文件入库与检索服务接口
type AgentServiceIFace interface {
	// Ingest 解析单个文件并入库（向量化后写入向量库）。
	// filename 用于识别文件类型，media 为文件原始字节。
	// 返回 ServiceResult：code=0 时 data 为 *IngestResult。
	Ingest(ctx context.Context, filename string, media []byte) (common.ServiceResult, error)

	// Search 根据自然语言查询检索相关文档片段。
	// topK<=0 时使用配置默认值。
	// 返回 ServiceResult：code=0 时 data 为 []SearchItem（按相关度降序）。
	Search(ctx context.Context, query string, topK int) (common.ServiceResult, error)
}

// IngestResult 入库结果
type IngestResult struct {
	// DocumentCount 入库的文件数量（当前接口单文件入库，恒为 1）
	DocumentCount int `json:"document_count"`
	// ChunkCount 切分并写入向量库的文本块数量
	ChunkCount int `json:"chunk_count"`
}

// SearchItem 检索单条结果
type SearchItem struct {
	// Content 命中的文本块内容
	Content string `json:"content"`
	// Score 相似度得分（越高越相关）
	Score float64 `json:"score"`
	// Source 来源文件名（入库时记录）
	Source string `json:"source"`
}
