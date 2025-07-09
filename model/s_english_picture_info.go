package model

// 英语绘本详情
type SEnglishPictureInfo struct {
	Id       int    `gorm:"column:id;type:int(11) unsigned;primary_key;AUTO_INCREMENT;comment:主键id" json:"id"`
	BookId   string `gorm:"column:book_id;type:char(32);comment:绘本id;NOT NULL" json:"book_id"`
	Pic      string `gorm:"column:pic;type:varchar(1024);comment:详情图(小);NOT NULL" json:"pic"`
	BPic     string `gorm:"column:b_pic;type:varchar(1024);comment:详情图(大);NOT NULL" json:"b_pic"`
	Mp3      string `gorm:"column:mp3;type:varchar(1024);comment:Mp3;NOT NULL" json:"mp3"`
	En       string `gorm:"column:en;type:text;comment:英文内容;NOT NULL" json:"en"`
	Zh       string `gorm:"column:zh;type:text;comment:中文内容;NOT NULL" json:"zh"`
	Position int    `gorm:"column:position;type:tinyint(4);default:0;comment:排序位置;NOT NULL" json:"position"`
}

func (m *SEnglishPictureInfo) TableName() string {
	return "s_english_picture_info"
}
