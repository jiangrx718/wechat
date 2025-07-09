package model

// 名称表
type SBookName struct {
	Id         uint   `gorm:"column:id;type:int(11) unsigned;primary_key;AUTO_INCREMENT;comment:主键id" json:"id"`
	CategoryId int    `gorm:"column:category_id;type:int(11);comment:类型" json:"category_id"`
	Name       string `gorm:"column:name;type:varchar(256);comment:名称" json:"name"`
	SSort      int    `gorm:"column:s_sort;type:tinyint(4);default:0;comment:排序" json:"s_sort"`
	SType      int    `gorm:"column:s_type;type:tinyint(4);default:0;comment:1中文绘本,2英文绘本,3古诗绘本" json:"s_type"`
	Status     int    `gorm:"column:status;type:tinyint(1);default:1;comment:状态" json:"status"`
}

func (m *SBookName) TableName() string {
	return "s_book_name"
}
