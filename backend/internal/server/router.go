// router.go 注册所有 API 路由分组，挂载中间件，统一管理 /api/v1 下的接口
package server

import (
	"ElainaWeb/global"

	"github.com/gin-gonic/gin"
)

func RouterInit() *gin.Engine {
	gin.SetMode(global.Config.System.Env)
	Router := gin.Default()

	return Router
}
