package response

import "github.com/jiangrx816/wechat/model"

type ChineseBookNavNameResponse struct {
	List []model.SBookName `json:"list"`
}

type ChineseBookResponse struct {
	Page      int                   `json:"page"`
	Total     int64                 `json:"total"`
	TotalPage int                   `json:"total_page"`
	List      []ResponseChineseBook `json:"list"`
}
type ResponseBookInfoCount struct {
	BookId    string `json:"book_id"`
	BookCount string `json:"book_count"`
}

type ResponseChineseBook struct {
	Id        int    `json:"-"`
	BookId    string `json:"book_id"`
	Title     string `json:"title"`
	Icon      string `json:"icon"`
	Level     int8   `json:"-"`
	Position  int    `json:"-"`
	BookCount string `json:"book_count"`
}

type ChineseBookInfoResponse struct {
	Info []ResponseChineseBookInfo `json:"info"`
}
type ResponseChineseBookInfo struct {
	Id       int    `json:"id"`
	BookId   string `json:"book_id"`
	Mp3      string `json:"mp3"`
	Pic      string `json:"pic"`
	Position uint8  `json:"position"`
}
