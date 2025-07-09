package wechat_service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/jiangrx816/wechat/common/response"
	"github.com/jiangrx816/wechat/core/server/api"
	"github.com/jiangrx816/wechat/utils"
	"github.com/jiangrx816/wechat/utils/errs"
)

func (ps *WechatService) FindPoetryBookList(page, limit, typeId int) (resp response.PoetryBookResponse, apiErr api.Error) {
	utils.DefaultIntOne(&page)
	utils.DefaultIntFifty(&limit)
	utils.DefaultIntOne(&typeId)
	// 设置表单数据
	data := url.Values{
		"org_id":  {"0"},
		"user_id": {"1"},
		"diff":    {"all"},
		"is_read": {"all"},
		"sort":    {"id"},
		"limit":   {"" + strconv.Itoa(limit) + ""},
		"page":    {"" + strconv.Itoa(page) + ""},
		"type_id": {"" + strconv.Itoa(typeId) + ""},
	}

	// 将表单数据转换为字符串
	dataString := data.Encode()

	// 转换为字节流
	dataBytes := strings.NewReader(dataString)

	// 设置请求
	req, err := http.NewRequest("POST", "https://mzbook.com/api/book.Book/getBookList", dataBytes)
	if err != nil {
		panic(err)
	}

	// 设置Content-Type
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 发送请求
	resp1, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp1.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(resp1.Body)
	if err != nil {
		panic(err)
	}

	// 打印响应
	var contentResult response.ResponsePoetryBookJson
	if err := json.Unmarshal(body, &contentResult); err != nil {
		return
	}
	jsonFilePath := fmt.Sprintf("/Users/jiang/小程序/poetry/info/%d.json", typeId)
	// 将 拓扑数据类型数据的 JSON 内容写入文件
	if err = os.WriteFile(jsonFilePath, body, 0644); err != nil {
		return resp, errs.NewError(err.Error())
	}

	resp.Data = contentResult.Data
	return
}

func (ps *WechatService) FindPoetryBookInfo(id int) (resp response.PoetryBookInfoResponse, apiErr api.Error) {
	utils.DefaultIntOne(&id)
	// 设置表单数据
	data := url.Values{
		"id": {"" + strconv.Itoa(id) + ""},
	}

	// 将表单数据转换为字符串
	dataString := data.Encode()

	// 转换为字节流
	dataBytes := strings.NewReader(dataString)

	// 设置请求
	req, err := http.NewRequest("POST", "https://mzbook.com/api/book.Book/getSentence", dataBytes)
	if err != nil {
		panic(err)
	}

	// 设置Content-Type
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 发送请求
	resp1, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp1.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(resp1.Body)
	if err != nil {
		panic(err)
	}

	// 打印响应
	var contentResult response.ResponsePoetryBookInfo
	if err := json.Unmarshal(body, &contentResult); err != nil {
		return
	}

	resp.Data = contentResult.Data
	return
}
