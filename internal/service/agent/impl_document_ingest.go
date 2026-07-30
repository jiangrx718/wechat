package agent

import (
	"bytes"
	"context"

	"wechat-tools/internal/common"
	"wechat-tools/internal/service/agent/loader"
	"wechat-tools/utils"

	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
)

// Ingest 解析单个文件、切分为文本块后写入向量库
func (s *AgentService) Ingest(ctx context.Context, filename string, media []byte) (common.ServiceResult, error) {
	var (
		logger = utils.SugarContext(ctx)
		result = common.NewServiceResult()
	)

	if s.InitErr() != nil {
		result.SetError(&common.ServiceError{Code: 500, Message: "知识库服务未就绪"}, s.InitErr())
		return result, nil
	}

	if len(media) == 0 {
		result.SetError(&common.ServiceError{Code: 400, Message: "文件内容为空"})
		return result, nil
	}

	// 1. 解析文件为纯文本
	text, err := loader.Load(filename, bytes.NewReader(media))
	if err != nil {
		logger.Errorw("Agent Ingest loader.Load error", "filename", filename, "error", err)
		result.SetError(&common.ServiceError{Code: 400, Message: "文件解析失败"}, err)
		return result, nil
	}
	if text == "" {
		result.SetError(&common.ServiceError{Code: 400, Message: "文件无可提取的文本内容"})
		return result, nil
	}

	// 2. 文本切块
	chunks := splitText(text, s.chunkSize, s.chunkOverlap)
	if len(chunks) == 0 {
		result.SetError(&common.ServiceError{Code: 400, Message: "文本切块为空"})
		return result, nil
	}

	// 3. 构造文档：每块对应一个 schema.Document
	docs := make([]*schema.Document, 0, len(chunks))
	for i, chunk := range chunks {
		docs = append(docs, &schema.Document{
			ID:      uuid.NewString(),
			Content: chunk,
			MetaData: map[string]any{
				"source": filename,
				"index":  i,
			},
		})
	}

	// 4. 写入向量库（indexer 内部会调用 embedder 向量化）
	ids, err := s.indexer.Store(ctx, docs)
	if err != nil {
		logger.Errorw("Agent Ingest indexer.Store error", "filename", filename, "error", err)
		result.SetError(&common.ServiceError{Code: 500, Message: "入库失败"}, err)
		return result, nil
	}

	logger.Infow("Agent Ingest success", "filename", filename, "chunk_count", len(ids))

	result.Data = &IngestResult{
		DocumentCount: 1,
		ChunkCount:    len(ids),
	}
	result.SetMessage("操作成功")
	return result, nil
}

// splitText 按字符数（rune）滑窗切分文本，相邻块间保留 overlap 个字符重叠，中文友好。
func splitText(text string, size, overlap int) []string {
	runes := []rune(text)
	if size <= 0 {
		return []string{text}
	}
	if len(runes) <= size {
		return []string{text}
	}

	step := size - overlap
	if step <= 0 {
		step = size
	}

	var chunks []string
	for i := 0; i < len(runes); i += step {
		end := i + size
		if end > len(runes) {
			end = len(runes)
		}
		chunks = append(chunks, string(runes[i:end]))
		if end == len(runes) {
			break
		}
	}

	return chunks
}
