package response

type IndexResponse struct {
	Total int64 `json:"total"`
	Data  any   `json:"data"`
}
