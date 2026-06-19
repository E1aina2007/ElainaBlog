package util

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var httpClient = &http.Client{
	Timeout: 5 * time.Second,
}

// 匹配 <link rel="icon" href="..."> 或 <link rel="shortcut icon" href="...">
var linkRe = regexp.MustCompile(`(?i)<link[^>]+rel=["'](?:shortcut )?icon["'][^>]+href=["']([^"']+)["']`)

// 匹配 <link href="..." rel="icon"> （href 在 rel 前面的情况）
var linkRe2 = regexp.MustCompile(`(?i)<link[^>]+href=["']([^"']+)["'][^>]+rel=["'](?:shortcut )?icon["']`)

// 匹配 <meta property="og:image" content="...">
var ogImageRe = regexp.MustCompile(`(?i)<meta[^>]+property=["']og:image["'][^>]+content=["']([^"']+)["']`)

// FaviconProxy 处理 GET /api/ui/favicon?domain=example.com
func FaviconProxy(c *gin.Context) {
	domain := c.Query("domain")
	if domain == "" {
		c.String(http.StatusBadRequest, "missing domain parameter")
		return
	}

	iconURL, err := FetchIconURL(c.Request.Context(), domain)
	if err != nil {
		c.String(http.StatusNotFound, "favicon not found")
		return
	}

	body, contentType, err := DownloadIcon(c.Request.Context(), iconURL, domain)
	if err != nil {
		c.String(http.StatusNotFound, "favicon not found")
		return
	}
	defer body.Close()

	c.Header("Content-Type", contentType)
	c.Header("Cache-Control", "public, max-age=86400")
	io.Copy(c.Writer, body)
}

// FetchIconURL 根据域名获取最佳 favicon URL，按优先级尝试多种策略
func FetchIconURL(ctx context.Context, domain string) (string, error) {
	// 策略1：解析 HTML 中的 <link rel="icon"> 标签（最准确）
	if iconURL, err := extractIconFromHTML(ctx, domain); err == nil && iconURL != "" {
		return iconURL, nil
	}

	// 策略2：直接请求常见路径
	paths := []string{"/favicon.ico", "/apple-touch-icon.png"}
	for _, path := range paths {
		iconURL := fmt.Sprintf("https://%s%s", domain, path)
		req, err := http.NewRequestWithContext(ctx, http.MethodHead, iconURL, nil)
		if err != nil {
			continue
		}
		req.Header.Set("Referer", fmt.Sprintf("https://%s/", domain))
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
		resp, err := httpClient.Do(req)
		if err != nil {
			continue
		}
		resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			return iconURL, nil
		}
	}

	return "", fmt.Errorf("favicon not found for %s", domain)
}

// DownloadIcon 下载图标并返回 body 和 Content-Type
func DownloadIcon(ctx context.Context, iconURL, domain string) (io.ReadCloser, string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, iconURL, nil)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("Referer", fmt.Sprintf("https://%s/", domain))
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, "", err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, "", fmt.Errorf("status %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/x-icon"
	}
	return resp.Body, contentType, nil
}

// extractIconFromHTML 请求首页 HTML，提取 <link rel="icon"> 中的图标 URL
func extractIconFromHTML(ctx context.Context, domain string) (string, error) {
	homeURL := fmt.Sprintf("https://%s/", domain)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, homeURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status %d", resp.StatusCode)
	}

	// 只读前 64KB，避免大页面浪费内存
	body, err := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
	if err != nil {
		return "", err
	}
	html := string(body)

	// 尝试 <link rel="icon">
	if m := linkRe.FindStringSubmatch(html); len(m) > 1 {
		return resolveURL(domain, m[1]), nil
	}
	// 尝试 href 在 rel 前面的情况
	if m := linkRe2.FindStringSubmatch(html); len(m) > 1 {
		return resolveURL(domain, m[1]), nil
	}
	// 回退：og:image
	if m := ogImageRe.FindStringSubmatch(html); len(m) > 1 {
		return resolveURL(domain, m[1]), nil
	}

	return "", nil
}

// resolveURL 将相对路径转为绝对 URL
func resolveURL(domain, href string) string {
	if strings.HasPrefix(href, "http://") || strings.HasPrefix(href, "https://") {
		return href
	}
	if strings.HasPrefix(href, "//") {
		return "https:" + href
	}
	if strings.HasPrefix(href, "/") {
		return fmt.Sprintf("https://%s%s", domain, href)
	}
	return fmt.Sprintf("https://%s/%s", domain, href)
}
