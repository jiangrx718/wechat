package request

// IndexPosttRequest 请求参数
type IndexPosttRequest struct {
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
	PlanId    string `json:"plan_id" binding:"required"`
	DataType  int    `json:"data_type" binding:"required"`
	Remark    string `json:"remark"`
}
