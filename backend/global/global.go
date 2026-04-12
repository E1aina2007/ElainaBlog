// global.go 集中管理全局变量
package global

import (
	"ElainaWeb/config"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Config *config.Config // 全局配置实例
	DB     *gorm.DB       // 全局数据库连接
	Logger *zap.Logger    // 全局日志实例
	Redis  *redis.Client  // 全局Redis连接
)
