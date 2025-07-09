package wechat_service

import (
	"math"

	"github.com/gin-gonic/gin"
	"github.com/jiangrx816/wechat/common"
	"github.com/jiangrx816/wechat/common/response"
	"github.com/jiangrx816/wechat/core/server/api"
	"github.com/jiangrx816/wechat/utils"
	"github.com/jiangrx816/wechat/utils/errs"
)

/**
 * @Description 获取初始的栏目列表
 */
func (ps *WechatService) ApiServiceGetCategoryList(ctx *gin.Context, typeId int) (resp response.ChineseBookNavNameResponse, apiErr api.Error) {
	utils.DefaultIntOne(&typeId)
	bookNameList, err := ps.ServiceDBFindCategoryList(typeId)
	if condition := err != nil; condition {
		apiErr = errs.NewError(err.Error())
		return
	}
	resp.List = bookNameList
	return
}

/**
 * @Description 获取中文绘本列表
 **/
func (ps *WechatService) ApiServiceChineseBookList(ctx *gin.Context, page, level int) (resp response.ChineseBookResponse, apiErr api.Error) {
	utils.DefaultIntOne(&page)
	utils.DefaultIntOne(&level)
	size := common.DEFAULT_PAGE_SIZE
	offset := size * (page - 1)

	// 从DB获取绘本列表数据
	total, bookList, err := ps.ServiceDBFindBookList(level, size, offset)
	if condition := err != nil; condition {
		apiErr = errs.NewError(err.Error())
		return
	}

	// 返回数据
	resp.Page = page
	resp.Total = total
	resp.List = bookList
	resp.TotalPage = int(math.Ceil(float64(total) / float64(size)))

	return
}

/**
 * @Description 获取中文绘本详情数据
 */
func (ps *WechatService) ApiServiceChineseBookInfo(ctx *gin.Context, bookId string) (resp response.ChineseBookInfoResponse, apiErr api.Error) {

	// 从DB获取绘本详情数据
	bookInfo, err := ps.ServiceDBFindBookInfo(bookId)

	if condition := err != nil; condition {
		apiErr = errs.NewError(err.Error())
		return
	}

	// 返回数据
	resp.Info = bookInfo

	return
}

/**
 * @Description 获取中文绘本搜索数据
 */
func (ps *WechatService) ApiServiceChineseBookSearch(ctx *gin.Context, page int, value string) (resp response.ChineseBookResponse, apiErr api.Error) {
	utils.DefaultIntOne(&page)
	size := 100
	offset := size * (page - 1)

	// 从DB获取绘本搜索数据
	total, bookList, err := ps.ServiceDBFindBookSearch(value, size, offset)
	if condition := err != nil; condition {
		apiErr = errs.NewError(err.Error())
		return
	}

	// 返回数据
	resp.Page = page
	resp.Total = total
	resp.List = bookList
	resp.TotalPage = int(math.Ceil(float64(total) / float64(size)))

	return
}

/**
 * @Description 获取英文绘本列表
 **/
func (ps *WechatService) ApiServiceEnglishBookList(ctx *gin.Context, page, level int) (resp response.EnglishBookResponse, apiErr api.Error) {
	utils.DefaultIntOne(&page)
	utils.DefaultIntOne(&level)
	size := common.DEFAULT_PAGE_SIZE
	offset := size * (page - 1)

	// 从DB获取绘本列表数据
	total, bookList, err := ps.ServiceDBFindEnglishBookList(level, size, offset)
	if condition := err != nil; condition {
		apiErr = errs.NewError(err.Error())
		return
	}

	// 返回数据
	resp.Page = page
	resp.Total = total
	resp.List = bookList
	resp.TotalPage = int(math.Ceil(float64(total) / float64(size)))

	return
}

/**
 * @Description 获取英文文绘本详情数据
 */
func (ps *WechatService) ApiServiceEnglishBookInfo(ctx *gin.Context, bookId string) (resp response.EnglishBookInfoResponse, apiErr api.Error) {

	// 从DB获取绘本详情数据
	bookInfo, err := ps.ServiceDBFindEnglishBookInfo(bookId)

	if condition := err != nil; condition {
		apiErr = errs.NewError(err.Error())
		return
	}

	// 返回数据
	resp.Info = bookInfo

	return
}

/**
 * @Description 获取英文绘本搜索数据
 */
func (ps *WechatService) ApiServiceEnglishBookSearch(ctx *gin.Context, page int, value string) (resp response.EnglishBookResponse, apiErr api.Error) {
	utils.DefaultIntOne(&page)
	size := 100
	offset := size * (page - 1)

	// 从DB获取绘本搜索数据
	total, bookList, err := ps.ServiceDBFindEnglishBookSearch(value, size, offset)
	if condition := err != nil; condition {
		apiErr = errs.NewError(err.Error())
		return
	}

	// 返回数据
	resp.Page = page
	resp.Total = total
	resp.List = bookList
	resp.TotalPage = int(math.Ceil(float64(total) / float64(size)))

	return
}
