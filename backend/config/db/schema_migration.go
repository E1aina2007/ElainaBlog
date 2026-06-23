// schema_migration.go 定义迁移记录表的 GORM 模型
package db

import "time"

// SchemaMigration 对应 schema_migrations 表，记录已执行的迁移版本
type SchemaMigration struct {
	Version   string    `gorm:"primaryKey;column:version"`
	AppliedAt time.Time `gorm:"column:applied_at;autoCreateTime"`
}

func (SchemaMigration) TableName() string { return "schema_migrations" }
