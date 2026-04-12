// gorm.go 数据库模型初始化，执行 AutoMigrate 自动建表
package model

import (
	"ElainaWeb/global"
)

// AutoMigrate 自动迁移所有数据模型，在 db.Init() 之后调用
func AutoMigrate() {
	// TODO: 添加模型结构体后在此注册
	// global.DB.AutoMigrate(&database.User{}, &database.Article{}, ...)
	_ = global.DB
}
