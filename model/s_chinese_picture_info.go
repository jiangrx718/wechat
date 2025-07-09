package model

// 国学绘本详情
type SChinesePictureInfo struct {
	Id        uint   `gorm:"column:id;type:int(11) unsigned;primary_key;AUTO_INCREMENT;comment:主键id" json:"id"`
	BookId    string `gorm:"column:book_id;type:char(32);comment:绘本id;NOT NULL" json:"book_id"`
	BookIdOld uint   `gorm:"column:book_id_old;type:int(11) unsigned;default:0;comment:绘本id;NOT NULL" json:"book_id_old"`
	Pic       string `gorm:"column:pic;type:varchar(1024);comment:详情图;NOT NULL" json:"pic"`
	Position  int    `gorm:"column:position;type:tinyint(4);default:0;comment:排序位置;NOT NULL" json:"position"`
	Mp3       string `gorm:"column:mp3;type:varchar(1024);comment:Mp3;NOT NULL" json:"mp3"`
	Status    int    `gorm:"column:status;type:tinyint(1);default:1;comment:启用状态，1启用,0禁用;NOT NULL" json:"status"`
}

func (m *SChinesePictureInfo) TableName() string {
	return "s_chinese_picture_info"
}
