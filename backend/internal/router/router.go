// router.go 注册所有 API 路由分组，挂载中间件，统一管理 /api/v1 下的接口
package router

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func RouterInit(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "ElainaWeb",
			"msg":     "你好, 我是Elaina!",
			"time":    time.Now().Format("2006-01-02 15:04:05"),
		})
	})
}
