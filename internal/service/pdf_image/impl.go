package pdf_image

// Service 内部实现结构体
type Service struct{}

// NewPdfImageService 创建 PDF 转图片服务实例
func NewPdfImageService() *Service {
	return &Service{}
}

// Ensure Service implements ServiceIFace at compile time
var _ ServiceIFace = (*Service)(nil)
