package wechat_user

import (
	"context"
	"wechat-tools/internal/common"
	"wechat-tools/internal/dao"
	"wechat-tools/utils"
)

type SWechatUserExistResp struct {
	Count int64 `json:"count"`
}

func (s *Service) Exist(ctx context.Context, deviceId string) (common.ServiceResult, error) {
	var (
		logger = utils.SugarContext(ctx)
		result = common.NewServiceResult()
	)

	// 校验用户是否存在
	count, err := dao.SWechatUser.Where(dao.SWechatUser.DeviceId.Eq(deviceId)).Count()
	if err != nil {
		logger.Errorw("WechatUserService Exist Count error", "device_id", deviceId, "error", err)
		return result, err
	}

	result.Data = SWechatUserExistResp{
		Count: count,
	}
	result.SetMessage("操作成功")
	return result, nil
}
