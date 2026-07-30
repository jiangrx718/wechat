package miniwechat

import (
	"wechat-tools/internal/dao"
	"wechat-tools/utils"

	"gorm.io/gorm"
)

type CheckImageService struct {
	db *gorm.DB
}

func NewCheckImageService() *CheckImageService {
	s := &CheckImageService{db: utils.DB()}
	dao.SetDefault(utils.DB())
	return s
}
