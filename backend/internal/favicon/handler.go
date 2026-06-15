// favicon 代理接口：带 Referer 头请求目标网站图标，绕过防盗链
package favicon

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 候选路径，按优先级排列
var candidatePaths = []string{
	"/favicon.ico",
	"/apple-touch-icon.png",
}

var client = &http.Client{
	Timeout: 5 * time.Second,
}

// Proxy 处理 GET /api/ui/favicon?domain=example.com
// 依次尝试候选路径，返回第一个成功的图标
func Proxy(c *gin.Context) {
	domain := c.Query("domain")
	if domain == "" {
		c.String(http.StatusBadRequest, "missing domain parameter")
		return
	}

	for _, path := range candidatePaths {
		iconURL := fmt.Sprintf("https://%s%s", domain, path)
		req, err := http.NewRequestWithContext(c.Request.Context(), http.MethodGet, iconURL, nil)
		if err != nil {
			continue
		}
		// 伪装 Referer 绕防盗链
		req.Header.Set("Referer", fmt.Sprintf("https://%s/", domain))
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

		resp, err := client.Do(req)
		if err != nil {
			continue
		}

		if resp.StatusCode == http.StatusOK {
			contentType := resp.Header.Get("Content-Type")
			if contentType == "" {
				contentType = "image/x-icon"
			}
			c.Header("Content-Type", contentType)
			c.Header("Cache-Control", "public, max-age=86400")
			io.Copy(c.Writer, resp.Body)
			resp.Body.Close()
			return
		}
		resp.Body.Close()
	}

	// 所有候选都失败，返回 404
	c.String(http.StatusNotFound, "favicon not found")
}
