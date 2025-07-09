package wechat_service

import (
	"strings"

	"github.com/jiangrx816/wechat/common/response"
	"github.com/jiangrx816/wechat/model"
)

/**
 * @Description 获取英文绘本列表数据
 */
func (ps *WechatService) ServiceDBFindEnglishBookList(tType, size, offset int) (total int64, bookList []response.ResponseEnglishBook, err error) {
	var bookInitList []model.SEnglishPicture
	db := model.Default().Model(&model.SEnglishPicture{}).Debug()
	db = db.Where("type = ? and status = 1", tType).Count(&total)
	db = db.Order("position desc").Limit(size).Offset(offset)
	db.Find(&bookInitList)

	var bookInfoCountList []response.ResponseBookInfoCount
	db1 := model.Default().Model(&model.SEnglishPictureInfo{}).Debug()
	db1.Raw("SELECT book_id,count(id) as book_count FROM s_english_picture_info where status = 1 GROUP BY book_id").Scan(&bookInfoCountList)

	var temp response.ResponseEnglishBook
	for _, item := range bookInitList {
		temp.Id = item.Id
		temp.BookId = item.BookId
		temp.Title = item.Title
		temp.Icon = strings.ReplaceAll(item.Icon, "https", "http")
		temp.Type = item.Type
		temp.Position = item.Position
		bookList = append(bookList, temp)
	}
	for index, item := range bookList {
		for _, it := range bookInfoCountList {
			if item.BookId == it.BookId {
				bookList[index].BookCount = it.BookCount
			}
		}
	}
	return
}

/**
 * @Description 获取英文绘本详情数据
 */
func (ps *WechatService) ServiceDBFindEnglishBookInfo(bookId string) (bookInfoList []response.ResponseEnglishBookInfo, err error) {
	db := model.Default().Model(&model.SEnglishPictureInfo{}).Debug()
	db = db.Where("book_id = ? and status = 1", bookId).Order("position asc")
	err = db.Find(&bookInfoList).Error

	for index, item := range bookInfoList {
		bookInfoList[index].Mp3 = strings.ReplaceAll(item.Mp3, "https", "http")
		bookInfoList[index].Pic = strings.ReplaceAll(item.Pic, "https", "http")
		bookInfoList[index].BPic = strings.ReplaceAll(item.BPic, "https", "http")
	}
	return
}

/**
 * @Description 获取英文绘本搜索数据
 */
func (ps *WechatService) ServiceDBFindEnglishBookSearch(name string, size, offset int) (total int64, bookList []response.ResponseEnglishBook, err error) {
	var bookInitList []model.SEnglishPicture
	db := model.Default().Model(&model.SEnglishPicture{}).Debug()
	db = db.Where("status = 1 and title like ?", "%"+name+"%")
	db = db.Count(&total)
	db = db.Order("type asc,position desc").Limit(size).Offset(offset)
	db.Find(&bookInitList)

	var bookInfoCountList []response.ResponseBookInfoCount
	db1 := model.Default().Model(&model.SEnglishPictureInfo{}).Debug()
	db1.Raw("SELECT book_id,count(id) as book_count FROM s_english_picture_info GROUP BY book_id").Scan(&bookInfoCountList)

	var temp response.ResponseEnglishBook
	for _, item := range bookInitList {
		temp.Id = item.Id
		temp.BookId = item.BookId
		temp.Title = item.Title
		temp.Icon = strings.ReplaceAll(item.Icon, "https", "http")
		temp.Type = item.Type
		temp.Position = item.Position
		bookList = append(bookList, temp)
	}
	for index, item := range bookList {
		for _, it := range bookInfoCountList {
			if item.BookId == it.BookId {
				bookList[index].BookCount = it.BookCount
			}
		}
	}

	return
}
