package wechat_user

import (
	"context"
	"wechat-tools/internal/common"
	"wechat-tools/internal/dao"
	"wechat-tools/model"
	"wechat-tools/utils"
)

type SWechatUserResp struct {
	UserName string `json:"user_name"`
	DeviceId string `json:"device_id"`
	Score    int    `json:"score"`
}

func (s *Service) Create(ctx context.Context, userName, deviceId string, score int) (common.ServiceResult, error) {
	var (
		logger = utils.SugarContext(ctx)
		result = common.NewServiceResult()
	)

	wechatUserData := model.SWechatUser{
		UserName: userName,
		DeviceId: deviceId,
		Score:    score,
	}

	// 校验用户是否存在
	count, err := dao.SWechatUser.Where(dao.SWechatUser.UserName.Eq(userName)).Count()
	if err != nil {
		logger.Errorw("WechatUserService Create Count error", "user_name", userName, "error", err)
		return result, err
	}
	if count > 0 {
		result.SetError(&common.ServiceError{Code: 400, Message: "用户已存在,请更换用户名"})
		return result, nil
	}

	if err := dao.SWechatUser.Create(&wechatUserData); err != nil {
		logger.Errorw("WechatUserService Create dao.Create error", "error", err)
		return result, err
	}

	result.Data = SWechatUserResp{
		UserName: wechatUserData.UserName,
		DeviceId: wechatUserData.DeviceId,
		Score:    wechatUserData.Score,
	}
	result.SetMessage("操作成功")
	return result, nil
}
