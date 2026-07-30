package agent

import (
	"context"

	"wechat-tools/internal/common"
	"wechat-tools/utils"
)

// Search 根据自然语言查询检索相关文档片段
func (s *AgentService) Search(ctx context.Context, query string, topK int) (common.ServiceResult, error) {
	var (
		logger = utils.SugarContext(ctx)
		result = common.NewServiceResult()
	)

	if s.InitErr() != nil {
		result.SetError(&common.ServiceError{Code: 500, Message: "知识库服务未就绪"}, s.InitErr())
		return result, nil
	}

	if query == "" {
		result.SetError(&common.ServiceError{Code: 400, Message: "查询内容不能为空"})
		return result, nil
	}

	// topK<=0 时使用配置默认值
	if topK <= 0 {
		topK = s.defaultTopK
	}

	docs, err := s.retriever.Retrieve(ctx, query)
	if err != nil {
		logger.Errorw("Agent Search retriever.Retrieve error", "query", query, "error", err)
		result.SetError(&common.ServiceError{Code: 500, Message: "检索失败"}, err)
		return result, nil
	}

	items := make([]SearchItem, 0, len(docs))
	for _, doc := range docs {
		source, _ := doc.MetaData["source"].(string)
		// 兼容自定义 converter 未回填 source 的情况：尝试从 metadata 取字符串
		if source == "" {
			if v, ok := doc.MetaData[fieldSource].(string); ok {
				source = v
			}
		}
		items = append(items, SearchItem{
			Content: doc.Content,
			Score:   doc.Score(),
			Source:  source,
		})
	}

	result.Data = items
	result.SetMessage("操作成功")
	return result, nil
}
