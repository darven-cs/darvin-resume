package export

import (
	"context"
	"fmt"
	"os"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/cdproto/page"
)

// PDFOptions PDF 导出配置
type PDFOptions struct {
	PaperWidth  float64 // inches (A4: 8.27)
	PaperHeight float64 // inches (A4: 11.69)
	Scale       float64 // DPI / 72
	PrintBg     bool    // 打印背景色
}

// DefaultPDFOptions 返回默认 A4 PDF 选项
func DefaultPDFOptions() *PDFOptions {
	return &PDFOptions{
		PaperWidth:  8.27,  // A4: 210mm = 8.27 inches
		PaperHeight: 11.69, // A4: 297mm = 11.69 inches
		Scale:       1.0,   // DPI 96 / 72
		PrintBg:     true,
	}
}

// ExportPDFFromHTML 使用 Chromedp 无头浏览器从 HTML 内容导出 PDF
// 使用临时文件方案避免 data: URL 编码破坏 CSS（{} 被编码成 %7B%7D）
func ExportPDFFromHTML(ctx context.Context, htmlContent string, outputPath string, opts *PDFOptions) error {
	if opts == nil {
		opts = DefaultPDFOptions()
	}

	// 验证输出路径安全性（防止路径穿越）
	if err := validateOutputPath(outputPath); err != nil {
		return err
	}

	// 写入临时 HTML 文件（避免 data: URL 编码问题）
	tmpFile, err := os.CreateTemp("", "darvin-resume-export-*.html")
	if err != nil {
		return fmt.Errorf("create temp file failed: %w", err)
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpPath) // 使用后清理

	if err := os.WriteFile(tmpPath, []byte(htmlContent), 0644); err != nil {
		return fmt.Errorf("write temp HTML failed: %w", err)
	}

	// 使用 file URL 加载（无编码问题）
	fileURL := "file://" + tmpPath

	// 创建无头浏览器分配器
	allocCtx, cancel := chromedp.NewExecAllocator(ctx,
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-setuid-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
	)
	defer cancel()

	// 创建浏览器上下文
	browserCtx, cancel2 := chromedp.NewContext(allocCtx)
	defer cancel2()

	// 设置超时
	browserCtx, cancel3 := context.WithTimeout(browserCtx, 60*1000000000) // 60s
	defer cancel3()

	var pdfBuf []byte
	err = chromedp.Run(browserCtx,
		// 使用 file URL 加载 HTML（避免 data: URL 编码破坏 CSS）
		chromedp.Navigate(fileURL),
		chromedp.WaitReady("body"),

		// 生成 PDF
		chromedp.ActionFunc(func(ctx context.Context) error {
			p := page.PrintToPDF().
				WithPrintBackground(opts.PrintBg).
				WithLandscape(false).
				WithPaperWidth(opts.PaperWidth).
				WithPaperHeight(opts.PaperHeight).
				WithMarginTop(0).
				WithMarginBottom(0).
				WithMarginLeft(0).
				WithMarginRight(0).
				WithScale(opts.Scale).
				WithPreferCSSPageSize(true)

			buf, _, err := p.Do(ctx)
			if err != nil {
				return fmt.Errorf("PrintToPDF failed: %w", err)
			}
			pdfBuf = buf
			return nil
		}),
	)

	if err != nil {
		return fmt.Errorf("chromedp run failed: %w", err)
	}

	// 确保输出目录存在
	dir := ""
	if idx := lastIndexByte(outputPath, '/'); idx >= 0 {
		dir = outputPath[:idx]
	}
	if dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("create output directory failed: %w", err)
		}
	}

	// 写入文件
	if err := os.WriteFile(outputPath, pdfBuf, 0644); err != nil {
		return fmt.Errorf("write file failed: %w", err)
	}

	return nil
}

// validateOutputPath 验证输出路径安全性（防止路径穿越攻击）
func validateOutputPath(path string) error {
	// 禁止路径穿越
	if containsPathTraversal(path) {
		return fmt.Errorf("output path contains invalid traversal: %s", path)
	}
	return nil
}

// containsPathTraversal 检查路径是否包含 ".." 等路径穿越字符
func containsPathTraversal(path string) bool {
	// 简单检查：禁止明显的路径穿越模式
	cleaned := path
	for {
		old := cleaned
		cleaned = replaceAll(cleaned, "..", "")
		if cleaned == old {
			break
		}
	}
	// 如果清理前后不一致，说明包含穿越字符
	if cleaned != path {
		return true
	}
	return false
}

// replaceAll 简单的字符串替换
func replaceAll(s, old, new string) string {
	result := ""
	for {
		idx := indexOf(s, old)
		if idx < 0 {
			return result + s
		}
		result += s[:idx] + new
		s = s[idx+len(old):]
	}
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func lastIndexByte(s string, c byte) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == c {
			return i
		}
	}
	return -1
}
