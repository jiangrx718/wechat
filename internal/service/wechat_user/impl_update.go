package wechat_user

import (
	"context"
	"wechat-tools/internal/common"
	"wechat-tools/internal/dao"
	"wechat-tools/utils"
)

func (s *Service) Update(ctx context.Context, userName string, score int) (common.ServiceResult, error) {
	var (
		logger = utils.SugarContext(ctx)
		result = common.NewServiceResult()
	)

	wechatUserDao := dao.SWechatUser

	// 校验用户是否存在
	count, err := wechatUserDao.Where(wechatUserDao.UserName.Eq(userName)).Count()
	if err != nil {
		logger.Errorw("WechatUserService Update Count error", "user_name", userName, "error", err)
		return result, err
	}
	if count == 0 {
		result.SetError(&common.ServiceError{Code: 400, Message: "用户不存在"})
		return result, nil
	}

	updates := map[string]interface{}{
		"score": score,
	}

	if _, err := wechatUserDao.Where(wechatUserDao.UserName.Eq(userName)).Updates(updates); err != nil {
		logger.Errorw("WechatUserService Update Updates error", "user_name", userName, "error", err)
		return result, err
	}

	result.Data = SWechatUserResp{
		UserName: userName,
		Score:    score,
	}

	result.SetMessage("操作成功")
	return result, nil
}
