// mysql.go 初始化 MySQL 数据库连接，提供全局 *gorm.DB 实例，并自动执行数据库迁移
package db

import (
	"ElainaBlog/config"
	"ElainaBlog/pkg/zaplogger"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB // DB 全局数据库实例

// ConnectDB 连接数据库并初始化 *gorm.DB，不执行迁移
func ConnectDB(dbConfig *config.DbConfig) error {
	dsn := dbConfig.GetDSN()

	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: NewGormLogger(),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 表名不加 s 后缀
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return err
	}

	// 获取底层 *sql.DB 配置连接池
	sqlDB, err := gormDB.DB()
	if err != nil {
		return err
	}
	if err = sqlDB.Ping(); err != nil {
		return err
	}
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)

	DB = gormDB
	return nil
}

// InitDB 连接数据库并自动执行迁移（向后兼容）
func InitDB(dbConfig *config.DbConfig) error {
	if err := ConnectDB(dbConfig); err != nil {
		return err
	}

	if err := Migrate(); err != nil {
		return fmt.Errorf("数据库迁移失败: %w", err)
	}

	return nil
}

// Migrate 执行数据库向上迁移，自动应用所有未执行的迁移脚本
func Migrate() error {
	if DB == nil {
		return fmt.Errorf("数据库连接未初始化")
	}

	// 创建迁移记录表（如果不存在）
	if err := createMigrationTable(); err != nil {
		return fmt.Errorf("创建迁移记录表失败: %w", err)
	}

	// 获取迁移脚本目录
	migrationDir, err := getMigrationDir()
	if err != nil {
		return err
	}

	// 获取所有 .up.sql 文件
	upFiles, err := getUpMigrationFiles(migrationDir)
	if err != nil {
		return err
	}

	// 获取已执行的迁移版本
	executedVersions, err := getExecutedVersions()
	if err != nil {
		return fmt.Errorf("获取已执行版本失败: %w", err)
	}

	// 执行未完成的迁移
	var applied int
	for _, file := range upFiles {
		version := extractVersion(file)
		if _, executed := executedVersions[version]; executed {
			continue
		}

		if err := executeMigrationFile(filepath.Join(migrationDir, file), version); err != nil {
			return fmt.Errorf("执行迁移 %s 失败: %w", file, err)
		}
		zaplogger.Logger.Info("数据库迁移成功", zap.String("version", version))
		applied++
	}

	if applied == 0 {
		zaplogger.Logger.Info("数据库已是最新，无需迁移")
	} else {
		zaplogger.Logger.Info("数据库迁移全部完成", zap.Int("applied", applied))
	}

	return nil
}

// DownMigrate 执行数据库向下迁移，回滚所有已执行的迁移脚本
func DownMigrate() error {
	if DB == nil {
		return fmt.Errorf("数据库连接未初始化")
	}

	if err := createMigrationTable(); err != nil {
		return fmt.Errorf("创建迁移记录表失败: %w", err)
	}

	migrationDir, err := getMigrationDir()
	if err != nil {
		return err
	}

	// 获取已执行的迁移版本（按版本号倒序）
	executedVersions, err := getExecutedVersionsOrdered()
	if err != nil {
		return fmt.Errorf("获取已执行版本失败: %w", err)
	}

	var rolledBack int
	for _, version := range executedVersions {
		downFile := filepath.Join(migrationDir, version+".down.sql")
		if _, err := os.Stat(downFile); os.IsNotExist(err) {
			zaplogger.Logger.Warn("跳过回滚: 下迁文件不存在", zap.String("version", version))
			continue
		}

		if err := executeMigrationFile(downFile, version); err != nil {
			return fmt.Errorf("回滚 %s 失败: %w", version, err)
		}

		if err := DB.Delete(&SchemaMigration{}, "version = ?", version).Error; err != nil {
			return fmt.Errorf("删除迁移记录失败: %w", err)
		}

		zaplogger.Logger.Info("回滚成功", zap.String("version", version))
		rolledBack++
	}

	if rolledBack == 0 {
		zaplogger.Logger.Info("没有可回滚的迁移")
	} else {
		zaplogger.Logger.Info("数据库回滚全部完成", zap.Int("rolled_back", rolledBack))
	}

	return nil
}

// createMigrationTable 创建迁移记录表
func createMigrationTable() error {
	if DB.Migrator().HasTable(&SchemaMigration{}) {
		return nil
	}
	return DB.Migrator().CreateTable(&SchemaMigration{})
}

// getExecutedVersions 获取已执行的迁移版本
func getExecutedVersions() (map[string]bool, error) {
	var records []SchemaMigration
	if err := DB.Find(&records).Error; err != nil {
		return nil, err
	}
	versions := make(map[string]bool, len(records))
	for _, r := range records {
		versions[r.Version] = true
	}
	return versions, nil
}

// getExecutedVersionsOrdered 获取已执行的迁移版本（按版本号倒序）
func getExecutedVersionsOrdered() ([]string, error) {
	var records []SchemaMigration
	if err := DB.Order("version DESC").Find(&records).Error; err != nil {
		return nil, err
	}
	versions := make([]string, len(records))
	for i, r := range records {
		versions[i] = r.Version
	}
	return versions, nil
}

// executeMigrationFile 执行单个迁移文件
func executeMigrationFile(filePath, version string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("读取文件失败: %w", err)
	}

	sqlContent := strings.TrimSpace(string(content))
	if sqlContent == "" {
		return nil
	}

	return DB.Transaction(func(tx *gorm.DB) error {
		queries := splitSQL(sqlContent)
		for _, query := range queries {
			query = strings.TrimSpace(query)
			if query == "" {
				continue
			}
			if err := tx.Exec(query).Error; err != nil {
				return fmt.Errorf("执行SQL失败: %w\nSQL: %s", err, query)
			}
		}

		// 记录迁移版本（只对 up 文件）
		if strings.HasSuffix(filePath, ".up.sql") {
			if err := tx.Create(&SchemaMigration{Version: version}).Error; err != nil {
				return fmt.Errorf("记录迁移版本失败: %w", err)
			}
		}

		return nil
	})
}

// getMigrationDir 获取迁移脚本目录的绝对路径
func getMigrationDir() (string, error) {
	if dir := os.Getenv("MIGRATION_DIR"); dir != "" {
		absDir, err := filepath.Abs(dir)
		if err != nil {
			return "", fmt.Errorf("获取绝对路径失败: %w", err)
		}
		if _, err := os.Stat(absDir); os.IsNotExist(err) {
			return "", fmt.Errorf("迁移目录不存在: %s", absDir)
		}
		return absDir, nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("获取当前工作目录失败: %w", err)
	}

	migrationDir := filepath.Join(cwd, "config", "db", "SQLscript")
	absDir, err := filepath.Abs(migrationDir)
	if err != nil {
		return "", fmt.Errorf("获取绝对路径失败: %w", err)
	}

	if _, err := os.Stat(absDir); os.IsNotExist(err) {
		return "", fmt.Errorf("迁移目录不存在: %s", absDir)
	}

	return absDir, nil
}

// getUpMigrationFiles 获取所有 .up.sql 文件并排序
func getUpMigrationFiles(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var upFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".up.sql") {
			upFiles = append(upFiles, file.Name())
		}
	}

	sort.Strings(upFiles)
	return upFiles, nil
}

// extractVersion 从文件名中提取版本号
func extractVersion(filename string) string {
	parts := strings.Split(filename, ".")
	if len(parts) > 0 {
		return parts[0]
	}
	return filename
}

// splitSQL 按分号分割SQL语句（支持注释和引号）
func splitSQL(sql string) []string {
	var queries []string
	var current strings.Builder
	inSingleQuote := false
	inDoubleQuote := false
	inLineComment := false
	inBlockComment := false

	for i := 0; i < len(sql); i++ {
		ch := sql[i]

		if inBlockComment {
			if ch == '*' && i+1 < len(sql) && sql[i+1] == '/' {
				inBlockComment = false
				i++
			}
			continue
		}

		if inLineComment {
			if ch == '\n' {
				inLineComment = false
			} else {
				continue
			}
		}

		if !inSingleQuote && !inDoubleQuote {
			if ch == '-' && i+1 < len(sql) && sql[i+1] == '-' {
				inLineComment = true
				i++
				continue
			}
			if ch == '/' && i+1 < len(sql) && sql[i+1] == '*' {
				inBlockComment = true
				i++
				continue
			}
		}

		if ch == '\'' && !inDoubleQuote {
			if i+1 < len(sql) && sql[i+1] == '\'' {
				current.WriteByte(ch)
				current.WriteByte(sql[i+1])
				i++
				continue
			}
			inSingleQuote = !inSingleQuote
		} else if ch == '"' && !inSingleQuote {
			inDoubleQuote = !inDoubleQuote
		}

		if ch == ';' && !inSingleQuote && !inDoubleQuote {
			q := strings.TrimSpace(current.String())
			if q != "" {
				queries = append(queries, q)
			}
			current.Reset()
		} else {
			current.WriteByte(ch)
		}
	}

	q := strings.TrimSpace(current.String())
	if q != "" {
		queries = append(queries, q)
	}

	return queries
}
