package wechat_user

import (
	"wechat-tools/internal/dao"
	"wechat-tools/utils"

	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewWechatUserService() *Service {
	s := &Service{db: utils.DB()}
	dao.SetDefault(utils.DB())
	return s
}
