package pdf_image

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"

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

// tryRemoveWatermark 尝试基于 PDF 结构移除水印。
// pdfcpu 通过解析 PDF 内部对象识别并删除水印，识别不到时原样保留 PDF，不会误伤正文。
// 任何错误（含识别不到第三方水印、PDF 加密/损坏）均安全回退为返回原始 pdfData，
// 调用方因此永远可以继续后续渲染流程，不会因去水印失败而中断。
func tryRemoveWatermark(pdfData []byte) []byte {
	var buf bytes.Buffer
	reader := bytes.NewReader(pdfData)
	// selectedPages 传 nil 表示作用于所有页；conf 传 nil 使用默认配置
	if err := api.RemoveWatermarks(reader, &buf, nil, nil); err != nil {
		// 安全回退：保留原始 PDF，仅记录日志
		log.Printf("pdf_image: remove watermark failed, fallback to original PDF: %v", err)
		return pdfData
	}
	if buf.Len() == 0 {
		log.Printf("pdf_image: remove watermark produced empty output, fallback to original PDF")
		return pdfData
	}
	return buf.Bytes()
}

// Convert 实现 PDF 转图片核心逻辑。
// 全程使用系统临时目录，转换完成后将图片读入内存并立即清理临时目录，不在服务端持久化任何文件。
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

	// 渲染前去水印：基于 PDF 结构精准移除，识别不到水印时原样保留，不会误伤正文
	if params.RemoveWatermark {
		pdfData = tryRemoveWatermark(pdfData)
	}

	// 使用系统临时目录，函数返回后自动清理，不在服务端持久化
	workDir, err := os.MkdirTemp("", "pdf_image_*")
	if err != nil {
		result.SetCode(500)
		result.SetMessage("创建临时目录失败")
		return result, fmt.Errorf("MkdirTemp: %w", err)
	}
	defer os.RemoveAll(workDir)

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
	matches, err := filepath.Glob(filepath.Join(workDir, "page*"+ext))
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
		data, err := os.ReadFile(match)
		if err != nil {
			result.SetCode(500)
			result.SetMessage("读取转换后的图片失败")
			return result, fmt.Errorf("ReadFile %s: %w", match, err)
		}
		images = append(images, PageImage{
			Page: pageNumFromPath(match),
			Data: data,
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
