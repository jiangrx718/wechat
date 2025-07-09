package response

type EnglishBookResponse struct {
	Page      int                   `json:"page"`
	Total     int64                 `json:"total"`
	TotalPage int                   `json:"total_page"`
	List      []ResponseEnglishBook `json:"list"`
}

type ResponseEnglishBook struct {
	Id        int    `json:"-"`
	BookId    string `json:"book_id"`
	Title     string `json:"title"`
	Icon      string `json:"icon"`
	Type      int    `json:"-"`
	Position  int    `json:"-"`
	BookCount string `json:"book_count"`
}

type EnglishBookInfoResponse struct {
	Info []ResponseEnglishBookInfo `json:"info"`
}
type ResponseEnglishBookInfo struct {
	Id       int    `json:"id"`
	BookId   string `json:"book_id"`
	Mp3      string `json:"mp3"`
	Pic      string `json:"pic"`
	BPic     string `json:"b_pic"`
	Position uint8  `json:"position"`
}
