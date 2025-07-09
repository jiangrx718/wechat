package wechat_service

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

var (
	accessToken    string
	accessTokenExp time.Time
	jsapiTicket    string
	jsapiTicketExp time.Time
	mutex          sync.Mutex
	client         = resty.New()
)

type accessTokenResp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
}

type jsapiTicketResp struct {
	Ticket    string `json:"ticket"`
	ExpiresIn int    `json:"expires_in"`
	ErrCode   int    `json:"errcode"`
	ErrMsg    string `json:"errmsg"`
}

func getAccessToken() (string, error) {
	mutex.Lock()
	defer mutex.Unlock()

	if time.Now().Before(accessTokenExp) && accessToken != "" {
		return accessToken, nil
	}

	appId := viper.GetViper().GetString("mini.app_id")
	appSecret := viper.GetViper().GetString("mini.app_secret")
	resp, err := client.R().
		SetQueryParams(map[string]string{
			"grant_type": "client_credential",
			"appid":      appId,
			"secret":     appSecret,
		}).
		Get("https://api.weixin.qq.com/cgi-bin/token")

	if err != nil {
		return "", err
	}

	var result accessTokenResp
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return "", err
	}

	if result.ErrCode != 0 {
		return "", fmt.Errorf("微信接口错误: %d %s", result.ErrCode, result.ErrMsg)
	}

	accessToken = result.AccessToken
	accessTokenExp = time.Now().Add(time.Duration(result.ExpiresIn-200) * time.Second)

	return accessToken, nil
}

func (ps *WechatService) ApiServiceGetJsapiTicket() (string, error) {
	mutex.Lock()
	defer mutex.Unlock()

	if time.Now().Before(jsapiTicketExp) && jsapiTicket != "" {
		return jsapiTicket, nil
	}

	token, err := getAccessToken()
	if err != nil {
		return "", err
	}

	resp, err := client.R().
		SetQueryParams(map[string]string{
			"access_token": token,
			"type":         "jsapi",
		}).
		Get("https://api.weixin.qq.com/cgi-bin/ticket/getticket")

	if err != nil {
		return "", err
	}

	var result jsapiTicketResp
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return "", err
	}

	if result.ErrCode != 0 {
		return "", fmt.Errorf("微信接口错误: %d %s", result.ErrCode, result.ErrMsg)
	}

	jsapiTicket = result.Ticket
	jsapiTicketExp = time.Now().Add(time.Duration(result.ExpiresIn-200) * time.Second)

	return jsapiTicket, nil
}

func (ps *WechatService) ApiServiceGenerateNonceStr(length int) string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	sb := strings.Builder{}
	for i := 0; i < length; i++ {
		sb.WriteByte(chars[rand.Intn(len(chars))])
	}
	return sb.String()
}

func (ps *WechatService) ApiServiceSalcSignature(ticket, nonceStr string, timestamp int64, urlStr string) string {
	params := map[string]string{
		"jsapi_ticket": ticket,
		"noncestr":     nonceStr,
		"timestamp":    fmt.Sprintf("%d", timestamp),
		"url":          urlStr,
	}

	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var strBuilder strings.Builder
	for i, k := range keys {
		strBuilder.WriteString(k)
		strBuilder.WriteString("=")
		strBuilder.WriteString(params[k])
		if i < len(keys)-1 {
			strBuilder.WriteString("&")
		}
	}

	h := sha1.New()
	h.Write([]byte(strBuilder.String()))
	return hex.EncodeToString(h.Sum(nil))
}
