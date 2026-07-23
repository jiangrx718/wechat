package check_image

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"sync"
	"time"

	"wechat-tools/internal/common"
	"wechat-tools/utils"

	"github.com/spf13/viper"
)

const (
	// 微信开放接口
	urlAccessToken = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
	urlImgSecCheck = "https://api.weixin.qq.com/wxa/img_sec_check?access_token=%s"

	// 微信 errcode
	wxErrCodeOK           = 0
	wxErrCodeRiskyContent = 87014 // 内容含有违法违规内容

	// 提前刷新 access_token：微信有效期 7200s，留 200s 余量
	tokenRefreshLead = 200 * time.Second

	httpTimeout = 15 * time.Second
)

// tokenCache 进程内缓存 access_token，避免每次检测都重新换取
type tokenCache struct {
	mu        sync.Mutex
	token     string
	expiresAt time.Time
}

var globalTokenCache tokenCache

// wxTokenResp 微信 access_token 接口返回
type wxTokenResp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"` // 秒
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
}

// wxImgSecCheckResp 微信 img_sec_check 接口返回
type wxImgSecCheckResp struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// getAccessToken 获取（并缓存）微信 access_token
func getAccessToken(ctx context.Context) (string, error) {
	logger := utils.SugarContext(ctx)

	globalTokenCache.mu.Lock()
	defer globalTokenCache.mu.Unlock()

	// 未过期直接复用
	if globalTokenCache.token != "" && time.Now().Before(globalTokenCache.expiresAt) {
		return globalTokenCache.token, nil
	}

	appid := viper.GetString("wechat.appid")
	secret := viper.GetString("wechat.secret")
	if appid == "" || secret == "" || appid == "your-appid" {
		return "", fmt.Errorf("wechat appid/secret 未配置")
	}

	url := fmt.Sprintf(urlAccessToken, appid, secret)
	client := &http.Client{Timeout: httpTimeout}
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求 access_token 失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取 access_token 响应失败: %w", err)
	}

	var tr wxTokenResp
	if err := json.Unmarshal(body, &tr); err != nil {
		return "", fmt.Errorf("解析 access_token 响应失败: %w, body=%s", err, string(body))
	}
	if tr.ErrCode != 0 || tr.AccessToken == "" {
		logger.Errorw("getAccessToken wechat error", "errcode", tr.ErrCode, "errmsg", tr.ErrMsg)
		return "", fmt.Errorf("微信返回错误: errcode=%d errmsg=%s", tr.ErrCode, tr.ErrMsg)
	}

	globalTokenCache.token = tr.AccessToken
	globalTokenCache.expiresAt = time.Now().Add(time.Duration(tr.ExpiresIn)*time.Second - tokenRefreshLead)
	logger.Debugw("getAccessToken ok", "expires_in", tr.ExpiresIn)
	return globalTokenCache.token, nil
}

// invalidateToken access_token 失效时清缓存，下次重新换取
func invalidateToken() {
	globalTokenCache.mu.Lock()
	globalTokenCache.token = ""
	globalTokenCache.expiresAt = time.Time{}
	globalTokenCache.mu.Unlock()
}

// Check 对图片做内容安全检测
func (s *Service) Check(ctx context.Context, media []byte, filename string) (common.ServiceResult, error) {
	var (
		logger = utils.SugarContext(ctx)
		result = common.NewServiceResult()
	)

	if len(media) == 0 {
		result.SetError(&common.ServiceError{Code: 400, Message: "图片内容为空"})
		return result, nil
	}

	// 最多重试一次：遇到 access_token 失效时清缓存重来
	for attempt := 0; attempt < 2; attempt++ {
		token, err := getAccessToken(ctx)
		if err != nil {
			logger.Errorw("Check getAccessToken error", "error", err)
			result.SetError(&common.ServiceError{Code: 500, Message: "获取 access_token 失败"}, err)
			return result, nil
		}

		pass, retriable, err := callImgSecCheck(ctx, token, media, filename)
		if err == nil {
			// 检测成功，返回结果
			result.SetCode(0)
			result.SetMessage("操作成功")
			result.Data = &CheckResult{Pass: pass}
			if !pass {
				logger.Infow("Check image risky", "filename", filename)
			}
			return result, nil
		}

		// access_token 失效（40001/42001 等）：清缓存后重试一次
		if retriable && attempt == 0 {
			logger.Warnw("Check token invalid, retrying", "error", err)
			invalidateToken()
			continue
		}
		logger.Errorw("Check callImgSecCheck error", "error", err)
		result.SetError(&common.ServiceError{Code: 500, Message: "内容检测服务异常"}, err)
		return result, nil
	}

	// 不应到达
	result.SetError(&common.ServiceError{Code: 500, Message: "内容检测服务异常"})
	return result, nil
}

// callImgSecCheck 调用微信 img_sec_check
// 返回: pass 是否合规; retriable 是否可重试(token 失效); err 调用错误
func callImgSecCheck(ctx context.Context, token string, media []byte, filename string) (pass bool, retriable bool, err error) {
	if filename == "" {
		filename = "image.jpg"
	}

	// 组装 multipart：字段名必须是 media
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("media", filename)
	if err != nil {
		return false, false, fmt.Errorf("创建表单文件失败: %w", err)
	}
	if _, err = part.Write(media); err != nil {
		return false, false, fmt.Errorf("写入图片数据失败: %w", err)
	}
	if err = writer.Close(); err != nil {
		return false, false, fmt.Errorf("关闭 multipart writer 失败: %w", err)
	}

	url := fmt.Sprintf(urlImgSecCheck, token)
	client := &http.Client{Timeout: httpTimeout}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return false, false, fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		return false, false, fmt.Errorf("请求 img_sec_check 失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, false, fmt.Errorf("读取响应失败: %w", err)
	}

	var cr wxImgSecCheckResp
	if err := json.Unmarshal(respBody, &cr); err != nil {
		return false, false, fmt.Errorf("解析响应失败: %w, body=%s", err, string(respBody))
	}

	switch cr.ErrCode {
	case wxErrCodeOK:
		return true, false, nil
	case wxErrCodeRiskyContent:
		// 含违规内容：不是错误，是检测结果「不通过」
		return false, false, nil
	default:
		// 40001/42001: access_token 失效或过期，可重试
		if cr.ErrCode == 40001 || cr.ErrCode == 42001 {
			return false, true, fmt.Errorf("access_token 失效: errcode=%d errmsg=%s", cr.ErrCode, cr.ErrMsg)
		}
		return false, false, fmt.Errorf("微信返回错误: errcode=%d errmsg=%s", cr.ErrCode, cr.ErrMsg)
	}
}
