package agent

import (
	"wechat-tools/server/http/response"
	"wechat-tools/utils"

	"github.com/gin-gonic/gin"
)

// SearchReq 语义检索请求
type SearchReq struct {
	// Query 查询内容（自然语言）
	Query string `json:"query" binding:"required"`
	// TopK 返回结果数量，<=0 时使用配置默认值
	TopK int `json:"top_k"`
}

// Search 语义检索接口
//
// 请求: application/json  {"query":"问题","top_k":5}
// 响应: {"code":0,"msg":"操作成功","data":[{"content":"...","score":0.9,"source":"a.txt"}]}
func (h *AgentHandler) Search(ctx *gin.Context) {
	var (
		reqBody SearchReq
		logger  = utils.SugarContext(ctx)
	)
	if err := ctx.Bind(&reqBody); err != nil {
		logger.Infow("Handler Agent Search ctx.Bind err", "error", err)
		response.ParameterError(ctx, err)
		return
	}

	result, err := h.service.Search(ctx, reqBody.Query, reqBody.TopK)
	if err != nil {
		logger.Errorw("Handler Agent Search service.Search error", "error", err)
		response.InternalError(ctx)
		return
	}

	if result.GetCode() != 0 {
		response.Failed(ctx, result.GetCode(), result.GetMessage(), nil)
		return
	}

	response.Successful(ctx, result.GetData())
}
