package pdf_image

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"wechat-tools/internal/common"
)

func defaultConvertParams() *ConvertParams {
	return &ConvertParams{
		Format: "png",
		DPI:    150,
	}
}

// popplerBin 查找可用的 poppler 转换工具，优先 pdftocairo
func popplerBin() (string, string, error) {
	if bin, err := exec.LookPath("pdftocairo"); err == nil {
		return bin, "pdftocairo", nil
	}
	if bin, err := exec.LookPath("pdftoppm"); err == nil {
		return bin, "pdftoppm", nil
	}
	return "", "", fmt.Errorf("pdf to image tool not found, please install poppler-utils (pdftocairo/pdftoppm)")
}

func imageMime(format string) string {
	switch format {
	case "jpeg", "jpg":
		return "image/jpeg"
	case "png":
		return "image/png"
	default:
		return "image/png"
	}
}

func imageExt(format string) string {
	switch format {
	case "jpeg", "jpg":
		return ".jpg"
	default:
		return ".png"
	}
}

func pageNumFromPath(path string) int {
	base := filepath.Base(path)
	name := strings.TrimSuffix(base, filepath.Ext(base))
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

// Convert 实现 PDF 转图片核心逻辑。
// 当 params.OutputDir 非空时，图片直接写入该目录且不会被自动清理（由调用方管理生命周期）；
// 否则使用系统临时目录且会在函数返回后清理。
func (s *Service) Convert(ctx context.Context, pdfData []byte, filename string, opts ...ConvertOption) (common.ServiceResult, error) {
	result := common.NewServiceResult()

	bin, binName, err := popplerBin()
	if err != nil {
		result.SetCode(500)
		result.SetMessage(err.Error())
		return result, err
	}

	params := defaultConvertParams()
	for _, opt := range opts {
		opt(params)
	}

	format := strings.ToLower(params.Format)
	if format != "png" && format != "jpeg" && format != "jpg" {
		result.SetCode(400)
		result.SetMessage(fmt.Sprintf("不支持的图片格式: %s，仅支持 png/jpeg", params.Format))
		return result, nil
	}

	// 清晰度档位优先：若指定了 quality 则覆盖 DPI，否则沿用 DPI 参数
	if q := dpiFromQuality(params.Quality); q > 0 {
		params.DPI = q
	}

	if params.DPI < 72 || params.DPI > 600 {
		result.SetCode(400)
		result.SetMessage(fmt.Sprintf("不支持的 DPI: %d，有效范围 72-600", params.DPI))
		return result, nil
	}

	// Page 校验：0 表示转换全部页，>0 表示转换指定页（页码从 1 开始）
	if params.Page < 0 {
		result.SetCode(400)
		result.SetMessage(fmt.Sprintf("不支持的页码: %d，页码从 1 开始，0 表示转换全部页", params.Page))
		return result, nil
	}

	// 确定工作目录
	workDir := params.OutputDir
	if workDir == "" {
		workDir, err = os.MkdirTemp("", "pdf_image_*")
		if err != nil {
			result.SetCode(500)
			result.SetMessage("创建临时目录失败")
			return result, fmt.Errorf("MkdirTemp: %w", err)
		}
		defer os.RemoveAll(workDir)
	} else {
		if err := os.MkdirAll(workDir, 0755); err != nil {
			result.SetCode(500)
			result.SetMessage("创建输出目录失败")
			return result, fmt.Errorf("MkdirAll: %w", err)
		}
	}

	// 写入 PDF
	pdfPath := filepath.Join(workDir, "input.pdf")
	if err := os.WriteFile(pdfPath, pdfData, 0600); err != nil {
		result.SetCode(500)
		result.SetMessage("写入 PDF 临时文件失败")
		return result, fmt.Errorf("WriteFile: %w", err)
	}

	outputPrefix := filepath.Join(workDir, "page")
	args := []string{
		fmt.Sprintf("-%s", format),
		"-r", strconv.Itoa(params.DPI),
	}
	// 指定页码时，使用 -f / -l 限定只转换该页（pdftocairo 与 pdftoppm 页码均从 1 开始）
	if params.Page > 0 {
		args = append(args, "-f", strconv.Itoa(params.Page), "-l", strconv.Itoa(params.Page))
	}
	args = append(args, pdfPath, outputPrefix)
	cmd := exec.CommandContext(ctx, bin, args...)

	out, err := cmd.CombinedOutput()
	if err != nil {
		result.SetCode(500)
		result.SetMessage("PDF 转图片失败")
		return result, fmt.Errorf("%s error: %w, output: %s", binName, err, string(out))
	}

	ext := imageExt(format)
	matches, err := filepath.Glob(filepath.Join(workDir, "page*" + ext))
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

	sort.Slice(matches, func(i, j int) bool {
		return pageNumFromPath(matches[i]) < pageNumFromPath(matches[j])
	})

	mime := imageMime(format)
	images := make([]PageImage, 0, len(matches))
	for _, match := range matches {
		images = append(images, PageImage{
			Page: pageNumFromPath(match),
			Path: match,
			Mime: mime,
		})
	}

	result.SetCode(0)
	result.SetMessage("操作成功")
	result.Data = &ConvertResult{
		TotalPages: len(images),
		Images:     images,
	}

	return result, nil
}
