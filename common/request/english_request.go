package request

type RequestEnglishParam struct {
	HTs        string `json:"h_ts"`
	HM         int64  `json:"h_m"`
	Zone       int32  `json:"zone"`
	HLc        string `json:"h_lc"`
	Uid        int64  `json:"uid"`
	Token      string `json:"token"`
	HCn        string `json:"h_cn"`
	HDt        int32  `json:"h_dt"`
	Cate       int32  `json:"cate"`
	Atype      int32  `json:"atype"`
	Source     int32  `json:"source"`
	Offset     int32  `json:"offset"`
	Difficulty int32  `json:"difficulty"`
}

type RequestEnglishInfoParam struct {
	HTs    string `json:"h_ts"`
	HM     int64  `json:"h_m"`
	Zone   int32  `json:"zone"`
	HLc    string `json:"h_lc"`
	Uid    int64  `json:"uid"`
	Token  string `json:"token"`
	HCH    string `json:"h_ch"`
	HDt    int32  `json:"h_dt"`
	Cate   int32  `json:"cate"`
	Did    string `json:"did"`
	HDid   string `json:"h_did"`
	Atype  int32  `json:"atype"`
	BookId int    `json:"bookid"`
	Limit  int    `json:"limit"`
}

type EnglishHandleDataRequest struct {
	Level    int    `json:"level" binding:"required"`
	FilePath string `json:"file_path" binding:"required"`
}
