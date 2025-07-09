package response

type PoetryBookResponse struct {
	Data []ResponsePoetryBook `json:"data"`
}
type ResponsePoetryBook struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	TypeID     string `json:"type_id"`
	Img        string `json:"img"`
	Difficulty string `json:"difficulty"`
	Author     string `json:"author"`
	FreeRank   string `json:"free_rank"`
	IsRead     int    `json:"is_read"`
	IsFree     int    `json:"is_free"`
	IsVip      int    `json:"is_vip"`
}
type ResponsePoetryBookJson struct {
	Code     int                  `json:"code"`
	Msg      string               `json:"msg"`
	Data     []ResponsePoetryBook `json:"data"`
	TypeList struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Img  string `json:"img"`
	} `json:"type_list"`
}

type PoetryBookInfoResponse struct {
	Data map[string]PoetryBookInfo `json:"data"`
}

type ResponsePoetryBookInfo struct {
	Code int                       `json:"code"`
	Msg  string                    `json:"msg"`
	Data map[string]PoetryBookInfo `json:"data"`
}

type PoetryBookInfo struct {
	ID           int         `json:"id"`
	BookID       int         `json:"book_id"`
	CnContent    string      `json:"cn_content"`
	EnContent    string      `json:"en_content"`
	SpellContent string      `json:"spell_content"`
	CnArr        string      `json:"cn_arr"`
	EnArr        string      `json:"en_arr"`
	PlayTime     int         `json:"play_time"`
	EnPlayTime   int         `json:"en_play_time"`
	Img          string      `json:"img"`
	Meaning      string      `json:"meaning"`
	MeaningCnArr interface{} `json:"meaning_cn_arr"`
	MeaningEn    string      `json:"meaning_en"`
	MeaningEnArr interface{} `json:"meaning_en_arr"`
	Rank         int         `json:"rank"`
	IsShow       int         `json:"is_show"`
}

type AlbumPictureInfoResponse struct {
	Info ResponseAlbumBookInfo `json:"info"`
}

type AlbumBookIndexResponse struct {
	Page      int                      `json:"page"`
	Total     int64                    `json:"total"`
	TotalPage int                      `json:"total_page"`
	List      []ResponseAlbumBookIndex `json:"list"`
}

type ResponseAlbumBookIndex struct {
	Id       int    `json:"-"`
	BookId   string `json:"book_id"`
	Title    string `json:"title"`
	Icon     string `json:"icon"`
	Position uint8  `json:"position"`
}

type AlbumBookListResponse struct {
	Page      int                     `json:"page"`
	Total     int64                   `json:"total"`
	TotalPage int                     `json:"total_page"`
	List      []ResponseAlbumBookInfo `json:"list"`
}

type ResponseAlbumBookInfo struct {
	Id       int    `json:"id"`
	BookId   string `json:"book_id"`
	Mp3      string `json:"mp3"`
	Pic      string `json:"pic"`
	Position uint8  `json:"position"`
	Title    string `json:"title"`
	Desc     string `json:"desc"`
	Duration string `json:"duration"`
}
