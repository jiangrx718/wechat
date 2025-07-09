package wechat_handler

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/jiangrx816/wechat/common/request"
	"github.com/jiangrx816/wechat/utils/errs"
	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
)

/**
 * ApiGetCategoryList
 * @Description 获取初始的栏目列表
 * @param ctx
 **/
func (ph *WechatHandler) ApiGetCategoryList(ctx *gin.Context) {
	typeId, _ := strconv.Atoi(ctx.Query("type"))
	response, err := ph.service.ApiServiceGetCategoryList(ctx, typeId)
	if err != nil {
		ctx.JSON(errs.ErrResp(err))
		return
	}
	ctx.JSON(errs.SucResp(response))
}

/**
 * ApiChineseBookList
 * @Description 获取中文绘本列表
 * @param ctx
 **/
func (ph *WechatHandler) ApiChineseBookList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	level, _ := strconv.Atoi(ctx.Query("level"))

	response, err := ph.service.ApiServiceChineseBookList(ctx, page, level)
	if err != nil {
		ctx.JSON(errs.ErrResp(err))
		return
	}
	ctx.JSON(errs.SucResp(response))
}

/**
 * ApiChineseBookInfo
 * @Description 获取中文绘本详情
 * @param ctx
 **/
func (ph *WechatHandler) ApiChineseBookInfo(ctx *gin.Context) {
	bookId := ctx.Query("book_id")
	response, err := ph.service.ApiServiceChineseBookInfo(ctx, bookId)
	if err != nil {
		ctx.JSON(errs.ErrResp(err))
		return
	}
	ctx.JSON(errs.SucResp(response))
}

/**
 * ApiChineseBookSearch
 * @Description 获取中文绘本搜索列表
 * @param ctx
 **/
func (ph *WechatHandler) ApiChineseBookSearch(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	value := ctx.Query("value")
	response, err := ph.service.ApiServiceChineseBookSearch(ctx, page, value)
	if err != nil {
		ctx.JSON(errs.ErrResp(err))
		return
	}
	ctx.JSON(errs.SucResp(response))
}

/**
 * ApiEnglishBookList
 * @Description 获取中文绘本列表
 * @param ctx
 **/
func (ph *WechatHandler) ApiEnglishBookList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	level, _ := strconv.Atoi(ctx.Query("level"))

	response, err := ph.service.ApiServiceEnglishBookList(ctx, page, level)
	if err != nil {
		ctx.JSON(errs.ErrResp(err))
		return
	}
	ctx.JSON(errs.SucResp(response))
}

/**
 * ApiEnglishBookInfo
 * @Description 获取英文绘本详情
 * @param ctx
 **/
func (ph *WechatHandler) ApiEnglishBookInfo(ctx *gin.Context) {
	bookId := ctx.Query("book_id")
	response, err := ph.service.ApiServiceEnglishBookInfo(ctx, bookId)
	if err != nil {
		ctx.JSON(errs.ErrResp(err))
		return
	}
	ctx.JSON(errs.SucResp(response))
}

/**
 * ApiEnglishBookSearch
 * @Description 获取英文绘本搜索列表
 * @param ctx
 **/
func (ph *WechatHandler) ApiEnglishBookSearch(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	value := ctx.Query("value")
	response, err := ph.service.ApiServiceEnglishBookSearch(ctx, page, value)
	if err != nil {
		ctx.JSON(errs.ErrResp(err))
		return
	}
	ctx.JSON(errs.SucResp(response))
}

/**
 * ApiEnglishHandleData
 * @Description 处理数据
 * @param ctx
 **/
func (ph *WechatHandler) ApiEnglishHandleData(ctx *gin.Context) {
	var json request.EnglishHandleDataRequest
	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(errs.SucErrResp("", "参数错误"))
		return
	}
	response, err := ph.service.ApiServiceEnglishHandleData(ctx, json.Level, json.FilePath)
	if err != nil {
		ctx.JSON(errs.ErrResp(err))
		return
	}
	ctx.JSON(errs.SucResp(response))
}

/**
 * ApiPoetryBookList
 * @Description 古诗绘本-列表
 * @param ctx
 **/
func (ph *WechatHandler) ApiPoetryBookList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	limit, _ := strconv.Atoi(ctx.Query("limit"))
	typeId, _ := strconv.Atoi(ctx.Query("type_id"))
	response, err := ph.service.FindPoetryBookList(page, limit, typeId)
	if err != nil {
		ctx.JSON(errs.ErrResp(err))
		return
	}
	ctx.JSON(errs.SucResp(response))
}

/**
 * ApiPoetryBookInfo
 * @Description 古诗绘本-详情
 * @param ctx
 **/
func (ph *WechatHandler) ApiPoetryBookInfo(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Query("book_id"))
	response, err := ph.service.FindPoetryBookInfo(id)
	if err != nil {
		ctx.JSON(errs.ErrResp(err))
		return
	}
	ctx.JSON(errs.SucResp(response))
}

/**
 * ApiMathCalculationList
 * @Description 数学计算-列表
 * @param ctx
 **/
func (ph *WechatHandler) ApiMathCalculationList(ctx *gin.Context) {
	// 值, 几十之内
	value, _ := strconv.Atoi(ctx.Query("value"))
	// 限制, 每次生成多少个
	limit, _ := strconv.Atoi(ctx.Query("limit"))
	// 方式, 1:5以内的加法 2:10以内的减法 3:10以内加减 4:20以内加法(不进位) 5:20以内减法(不退位)
	// 6:20以内加法(进位) 7:20以内减法(退位) 8:20以内加减 9:100以内加法 10:100以内减法 11:100以内加减
	forward, _ := strconv.Atoi(ctx.Query("forward"))
	response, err := ph.service.ApiServiceMathCalculationList(ctx, forward, value, limit)
	if err != nil {
		ctx.JSON(errs.ErrResp(err))
		return
	}
	ctx.JSON(errs.SucResp(response))
}

/**
 * ApiWechatSignature
 * @Description: 微信签名
 * @param ctx
 */
func (ph *WechatHandler) ApiWechatSignature(ctx *gin.Context) {
	rawURL := ctx.Query("url")
	if rawURL == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "缺少url参数"})
		return
	}

	// url参数需要解码一次
	urlDecoded, err := url.QueryUnescape(rawURL)
	if err != nil {
		urlDecoded = rawURL // 失败则用原始
	}

	ticket, err := ph.service.ApiServiceGetJsapiTicket()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取jsapi_ticket失败：" + err.Error()})
		return
	}

	nonceStr := ph.service.ApiServiceGenerateNonceStr(16)
	timestamp := time.Now().Unix()
	signature := ph.service.ApiServiceSalcSignature(ticket, nonceStr, timestamp, urlDecoded)

	appId := viper.GetViper().GetString("mini.app_id")
	ctx.JSON(http.StatusOK, gin.H{
		"appId":     appId,
		"nonceStr":  nonceStr,
		"timestamp": timestamp,
		"signature": signature,
	})
}
