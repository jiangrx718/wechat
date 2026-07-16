package picture_book

import (
	"wechat-tools/internal/dao"
	"wechat-tools/utils"

	"gorm.io/gorm"
)

// Service 绘本服务
type Service struct {
	db *gorm.DB
}

// NewPictureBookService 创建绘本服务
func NewPictureBookService() *Service {
	s := &Service{db: utils.DB()}
	dao.SetDefault(utils.DB())
	return s
}
