package wechat_user

import (
	"context"
	"wechat-tools/internal/common"
	"wechat-tools/internal/dao"
	"wechat-tools/utils"
)

type SWechatUserUpdateResp struct {
	UserName string `json:"user_name"`
	Score    int    `json:"score"`
}

func (s *Service) Update(ctx context.Context, userName string, score int) (common.ServiceResult, error) {
	var (
		logger = utils.SugarContext(ctx)
		result = common.NewServiceResult()
	)

	wechatUserDao := dao.SWechatUser

	// 校验用户是否存在并获取当前数据
	user, err := wechatUserDao.Where(wechatUserDao.UserName.Eq(userName)).First()
	if err != nil {
		logger.Errorw("WechatUserService Update First error", "user_name", userName, "error", err)
		result.SetError(&common.ServiceError{Code: 400, Message: "用户不存在"})
		return result, nil
	}

	// 如果数据库中的score小于传递过来的score，直接返回
	if user.Score > score {
		result.SetError(&common.ServiceError{Code: 400, Message: "本局分值低于过往最高分，继续加油💪🏻"})
		return result, nil
	}

	updates := map[string]interface{}{
		"score": score,
	}

	if _, err := wechatUserDao.Where(wechatUserDao.UserName.Eq(userName)).Updates(updates); err != nil {
		logger.Errorw("WechatUserService Update Updates error", "user_name", userName, "error", err)
		return result, err
	}

	result.Data = SWechatUserUpdateResp{
		UserName: userName,
		Score:    score,
	}

	result.SetMessage("操作成功")
	return result, nil
}
