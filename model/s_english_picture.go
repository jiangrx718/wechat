package model

import (
	"time"
)

// 英语绘本
type SEnglishPicture struct {
	Id         int       `gorm:"column:id;type:int(11) unsigned;primary_key;AUTO_INCREMENT;comment:主键id" json:"id"`
	BookId     string    `gorm:"column:book_id;type:char(32);comment:绘本id;NOT NULL" json:"book_id"`
	Title      string    `gorm:"column:title;type:varchar(1024);comment:标题;NOT NULL" json:"title"`
	Icon       string    `gorm:"column:icon;type:varchar(1024);comment:封面图;NOT NULL" json:"icon"`
	Type       int       `gorm:"column:type;type:tinyint(1);default:1;comment:级别;NOT NULL" json:"type"`
	Position   int       `gorm:"column:position;type:int(11);default:0;comment:排序位置;NOT NULL" json:"position"`
	Status     int       `gorm:"column:status;type:tinyint(1);default:1;comment:状态:1启用,0禁用;NOT NULL" json:"status"`
	AddTime    time.Time `gorm:"column:add_time;type:timestamp;default:CURRENT_TIMESTAMP;comment:添加时间;NOT NULL" json:"add_time"`
	UpdateTime time.Time `gorm:"column:update_time;type:timestamp;default:CURRENT_TIMESTAMP;comment:更新时间;NOT NULL" json:"update_time"`
}

func (m *SEnglishPicture) TableName() string {
	return "s_english_picture"
}
