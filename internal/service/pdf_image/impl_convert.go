package pdf_image

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"

	"wechat-tools/internal/common"
)

// defaultConvertParams 返回默认转换参数
func defaultConvertParams() *ConvertParams {
	return &ConvertParams{
		Format: "png",
		DPI:    150,
	}
}

// popplerBin 查找可用的 poppler 转换工具，优先使用 pdftocairo
func popplerBin() (string, string, error) {
	if bin, err := exec.LookPath("pdftocairo"); err == nil {
		return bin, "pdftocairo", nil
	}
	if bin, err := exec.LookPath("pdftoppm"); err == nil {
		return bin, "pdftoppm", nil
	}
	return "", "", fmt.Errorf("pdf to image tool not found, please install poppler-utils (pdftocairo/pdftoppm)")
}

// imageMime 根据格式字符串返回对应的 MIME 类型
func imageMime(format string) string {
	switch format {
	case "jpeg", "jpg":
		return "image/jpeg"
	case "png":
		return "image/png"
	case "tiff":
		return "image/tiff"
	default:
		return "image/png"
	}
}

// imageExt 根据格式字符串返回文件扩展名
func imageExt(format string) string {
	switch format {
	case "jpeg", "jpg":
		return ".jpg"
	case "tiff":
		return ".tif"
	default:
		return ".png"
	}
}

// pageNumFromPath 从 pdftoppm/pdftocairo 输出的文件路径中提取页码。
// 输出格式为 <prefix>-<N>.<ext>，提取末尾的数字 N。
func pageNumFromPath(path string) int {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)
	parts := strings.Split(name, "-")
	if len(parts) == 0 {
		return 0
	}
	num, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		return 0
	}
	return num
}

// Convert 实现 PDF 转图片核心逻辑
func (s *Service) Convert(ctx context.Context, pdfData []byte, filename string, opts ...ConvertOption) (common.ServiceResult, error) {
	result := common.NewServiceResult()

	bin, binName, err := popplerBin()
	if err != nil {
		result.SetCode(500)
		result.SetMessage(err.Error())
		return result, err
	}

	// 合并参数选项
	params := defaultConvertParams()
	for _, opt := range opts {
		opt(params)
	}

	// 校验图片格式
	format := strings.ToLower(params.Format)
	if format != "png" && format != "jpeg" && format != "jpg" {
		result.SetCode(400)
		result.SetMessage(fmt.Sprintf("不支持的图片格式: %s，仅支持 png/jpeg", params.Format))
		return result, nil
	}

	// 校验 DPI 范围
	if params.DPI < 72 || params.DPI > 600 {
		result.SetCode(400)
		result.SetMessage(fmt.Sprintf("不支持的 DPI: %d，有效范围 72-600", params.DPI))
		return result, nil
	}

	// 创建临时工作目录
	workDir, err := os.MkdirTemp("", "pdf_image_*")
	if err != nil {
		result.SetCode(500)
		result.SetMessage("创建临时目录失败")
		return result, fmt.Errorf("MkdirTemp: %w", err)
	}
	defer os.RemoveAll(workDir)

	// 将 PDF 内容写入临时文件
	pdfPath := filepath.Join(workDir, "input.pdf")
	if err := os.WriteFile(pdfPath, pdfData, 0600); err != nil {
		result.SetCode(500)
		result.SetMessage("写入 PDF 临时文件失败")
		return result, fmt.Errorf("WriteFile: %w", err)
	}

	// 构建转换命令
	outputPrefix := filepath.Join(workDir, "page")
	args := []string{
		fmt.Sprintf("-%s", format),
		"-r", strconv.Itoa(params.DPI),
		pdfPath,
		outputPrefix,
	}
	cmd := exec.CommandContext(ctx, bin, args...)

	out, err := cmd.CombinedOutput()
	if err != nil {
		result.SetCode(500)
		result.SetMessage("PDF 转图片失败")
		return result, fmt.Errorf("%s error: %w, output: %s", binName, err, string(out))
	}

	// 收集生成的图片文件
	ext := imageExt(format)
	pattern := fmt.Sprintf("page*%s", ext)
	matches, err := filepath.Glob(filepath.Join(workDir, pattern))
	if err != nil {
		result.SetCode(500)
		result.SetMessage("读取转换结果失败")
		return result, fmt.Errorf("Glob: %w", err)
	}

	if len(matches) == 0 {
		result.SetCode(500)
		result.SetMessage("转换后未生成图片")
		return result, fmt.Errorf("no images generated from PDF")
	}

	// 按页码排序
	sort.Slice(matches, func(i, j int) bool {
		return pageNumFromPath(matches[i]) < pageNumFromPath(matches[j])
	})

	mime := imageMime(format)

	// 用 goroutine 并发读取并 base64 编码各个图片
	type readResult struct {
		data string
		err  error
	}
	results := make([]readResult, len(matches))
	sem := make(chan struct{}, runtime.NumCPU())
	var wg sync.WaitGroup

	for i, match := range matches {
		wg.Add(1)
		go func(idx int, path string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			raw, err := os.ReadFile(path)
			if err != nil {
				results[idx] = readResult{err: fmt.Errorf("ReadFile %s: %w", path, err)}
				return
			}
			results[idx] = readResult{data: base64.StdEncoding.EncodeToString(raw)}
		}(i, match)
	}
	wg.Wait()

	images := make([]PageImage, 0, len(matches))
	for i, r := range results {
		if r.err != nil {
			result.SetCode(500)
			result.SetMessage("读取图片文件失败")
			return result, r.err
		}
		images = append(images, PageImage{
			Page: pageNumFromPath(matches[i]),
			Data: r.data,
			Mime: mime,
		})
	}

	cr := &ConvertResult{
		TotalPages: len(images),
		Images:     images,
	}

	result.SetCode(0)
	result.SetMessage("操作成功")
	result.Data = cr
	return result, nil
}
