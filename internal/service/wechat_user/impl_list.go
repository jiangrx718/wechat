package wechat_user

import (
	"context"
	"wechat-tools/internal/common"
	"wechat-tools/internal/dao"
	"wechat-tools/model"
	"wechat-tools/utils"

	"gorm.io/gen"
)

type ListResponseData struct {
	List []WechatUserItem `json:"list"`
}

type WechatUserItem struct {
	UserName string `json:"user_name"`
	Score    int    `json:"score"`
}

func toWechatUserItem(m *model.SWechatUser) WechatUserItem {
	return WechatUserItem{
		UserName: m.UserName,
		Score:    m.Score,
	}
}

func (s *Service) List(ctx context.Context) (common.ServiceResult, error) {
	var (
		logger = utils.SugarContext(ctx)
		result = common.NewServiceResult()
	)

	wechatUserDao := dao.SWechatUser
	where := []gen.Condition{}

	list, _, err := wechatUserDao.Where(where...).Order(wechatUserDao.Score.Desc()).Debug().FindByPage(0, 10)
	if err != nil {
		logger.Errorw("WechatUserService List FindByPage error", "error", err)
		return result, err
	}

	items := make([]WechatUserItem, 0, len(list))
	for _, m := range list {
		items = append(items, toWechatUserItem(m))
	}

	result.Data = ListResponseData{
		List: items,
	}
	result.SetMessage("操作成功")
	return result, nil
}
