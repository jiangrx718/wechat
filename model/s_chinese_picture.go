package model

// 国学绘本
type SChinesePicture struct {
	Id        int    `gorm:"column:id;type:int(11) unsigned;primary_key;AUTO_INCREMENT;comment:主键id" json:"id"`
	BookId    string `gorm:"column:book_id;type:char(32);comment:绘本id;NOT NULL" json:"book_id"`
	BookIdOld int    `gorm:"column:book_id_old;type:int(11) unsigned;default:0;comment:绘本id;NOT NULL" json:"book_id_old"`
	Title     string `gorm:"column:title;type:varchar(1024);comment:标题;NOT NULL" json:"title"`
	Icon1     string `gorm:"column:icon_1;type:varchar(1024);comment:封面图;NOT NULL" json:"icon_1"`
	Icon      string `gorm:"column:icon;type:varchar(1024);comment:封面图;NOT NULL" json:"icon"`
	Type      int8   `gorm:"column:type;type:tinyint(1);default:1;comment:类型;NOT NULL" json:"type"`
	Position  int    `gorm:"column:position;type:int(11);default:0;comment:排序位置;NOT NULL" json:"position"`
	Status    int    `gorm:"column:status;type:tinyint(1);default:1;comment:启用状态，1启用,0禁用;NOT NULL" json:"status"`
}

func (m *SChinesePicture) TableName() string {
	return "s_chinese_picture"
}
