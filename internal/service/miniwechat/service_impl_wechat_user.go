package miniwechat

import (
	"wechat-tools/internal/dao"
	"wechat-tools/utils"

	"gorm.io/gorm"
)

type WechatUserService struct {
	db *gorm.DB
}

func NewWechatUserService() *WechatUserService {
	s := &WechatUserService{db: utils.DB()}
	dao.SetDefault(utils.DB())
	return s
}
