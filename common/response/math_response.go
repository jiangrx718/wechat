package response

type MathResponse struct {
	NumberOne int    `json:"number_one"`
	NumberTwo int    `json:"number_two"`
	Symbol    string `json:"symbol"`
	Result    int    `json:"result"`
}
