package wechat_service

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jiangrx816/wechat/core/server/api"
	"github.com/jiangrx816/wechat/model"
	"github.com/jiangrx816/wechat/utils"
)

/**
 * @Description 处理数据
 */
func (ps *WechatService) ApiServiceEnglishHandleData(ctx *gin.Context, level int, filePath string) (resp string, apiErr api.Error) {

	resp = "start english handle data"
	// 异步处理数据
	go func() {
		ps.HandleEnglishDataInfo(ctx, level, filePath)
	}()

	return
}

func (ps *WechatService) HandleEnglishDataInfo(ctx *gin.Context, level int, filePath string) (err error) {
	var bookInitList []model.SEnglishPictureInfo
	db := model.Default().Model(&model.SEnglishPictureInfo{}).Debug()
	db = db.Where("id >0").Order("id asc")
	db.Find(&bookInitList)

	for _, book := range bookInitList {
		dirPath := "/Users/jiang/小程序/file/" + book.BookId
		savePath := dirPath + "/" + strconv.Itoa(int(book.Id)) + "_b.jpg"
		flag, _ := utils.FileExists(savePath)
		if condition := !flag; condition {
			fmt.Printf("文件不存在: %s\n", savePath)
			ps.downloadAndSave(book)
		}
	}

	return nil
}

// 下载和保存 MP3 文件
func (ps *WechatService) downloadAndSave(book model.SEnglishPictureInfo) error {
	fileUrl := book.BPic
	dirPath := "/Users/jiang/小程序/file/" + book.BookId
	utils.ExistDir(dirPath)
	savePath := dirPath + "/" + strconv.Itoa(int(book.Id)) + "_b.jpg"

	// 创建 HTTP 客户端
	client := &http.Client{}
	resp, err := client.Get(fileUrl)
	if err != nil {
		return fmt.Errorf("下载请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("服务器返回错误状态码: %d", resp.StatusCode)
	}

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(savePath), 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	// 创建文件
	outFile, err := os.Create(savePath)
	if err != nil {
		return fmt.Errorf("文件创建失败: %w", err)
	}
	defer outFile.Close()

	// 写入内容
	if _, err := io.Copy(outFile, resp.Body); err != nil {
		return fmt.Errorf("文件写入失败: %w", err)
	}

	return nil
}
