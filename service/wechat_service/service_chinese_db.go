package wechat_service

import (
	"github.com/jiangrx816/wechat/common/response"
	"github.com/jiangrx816/wechat/model"
)

/**
 * @Description 获取初始的栏目列表
 */
func (ps *WechatService) ServiceDBFindCategoryList(typeId int) (bookNameList []model.SBookName, err error) {
	db := model.Default().Model(&model.SBookName{}).Debug()
	db = db.Where("status = 1 and s_type = ?", typeId)
	db = db.Order("s_sort asc").Order("id asc")
	err = db.Find(&bookNameList).Error
	return
}

/**
 * @Description 获取中文绘本列表数据
 */
func (ps *WechatService) ServiceDBFindBookList(tType, size, offset int) (total int64, bookList []response.ResponseChineseBook, err error) {
	var bookInitList []model.SChinesePicture
	db := model.Default().Model(&model.SChinesePicture{}).Debug()
	db = db.Where("type = ? and status = 1", tType).Count(&total)
	db = db.Order("position desc").Limit(size).Offset(offset)
	db.Find(&bookInitList)

	var bookInfoCountList []response.ResponseBookInfoCount
	db1 := model.Default().Model(&model.SChinesePictureInfo{}).Debug()
	db1.Raw("SELECT book_id,count(id) as book_count FROM s_chinese_picture_info where status = 1 GROUP BY book_id").Scan(&bookInfoCountList)

	var temp response.ResponseChineseBook
	for _, item := range bookInitList {
		temp.Id = item.Id
		temp.BookId = item.BookId
		temp.Title = item.Title
		temp.Icon = item.Icon
		temp.Level = item.Type
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
 * @Description 获取中文绘本详情数据
 */
func (ps *WechatService) ServiceDBFindBookInfo(bookId string) (bookInfoList []response.ResponseChineseBookInfo, err error) {
	db := model.Default().Model(&model.SChinesePictureInfo{}).Debug()
	db = db.Where("book_id = ? and status = 1", bookId).Order("position asc")
	err = db.Find(&bookInfoList).Error
	return
}

/**
 * @Description 获取中文绘本搜索数据
 */
func (ps *WechatService) ServiceDBFindBookSearch(name string, size, offset int) (total int64, bookList []response.ResponseChineseBook, err error) {
	var bookInitList []model.SChinesePicture
	db := model.Default().Model(&model.SChinesePicture{}).Debug()
	db = db.Where("status = 1 and title like ?", "%"+name+"%")
	db = db.Count(&total)
	db = db.Order("type asc,position desc").Limit(size).Offset(offset)
	db.Find(&bookInitList)

	var bookInfoCountList []response.ResponseBookInfoCount
	db1 := model.Default().Model(&model.SChinesePictureInfo{}).Debug()
	db1.Raw("SELECT book_id,count(id) as book_count FROM s_chinese_picture_info GROUP BY book_id").Scan(&bookInfoCountList)

	var temp response.ResponseChineseBook
	for _, item := range bookInitList {
		temp.Id = item.Id
		temp.BookId = item.BookId
		temp.Title = item.Title
		temp.Icon = item.Icon
		temp.Level = item.Type
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
