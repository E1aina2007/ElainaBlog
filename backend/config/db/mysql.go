// mysql.go 初始化 MySQL 数据库连接，提供全局 *sql.DB 实例，并自动执行数据库迁移
package db

import (
	"ElainaBlog/config"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var DBPool *sql.DB // DBPool 全局数据库连接池

func InitDB(dbConfig *config.DbConfig) error {
	dsn := dbConfig.GetDSN()
	db, err := sql.Open(dbConfig.SqlName, dsn)
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	db.SetMaxIdleConns(dbConfig.MaxIdleConns)
	db.SetMaxOpenConns(dbConfig.MaxOpenConns)

	DBPool = db

	// 自动执行数据库迁移
	if err := Migrate(); err != nil {
		return fmt.Errorf("数据库迁移失败: %w", err)
	}

	return nil
}

// Migrate 执行数据库向上迁移，自动应用所有未执行的迁移脚本
func Migrate() error {
	if DBPool == nil {
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
	for _, file := range upFiles {
		version := extractVersion(file)
		if _, executed := executedVersions[version]; executed {
			continue
		}

		fmt.Printf("执行迁移: %s\n", file)
		if err := executeMigrationFile(filepath.Join(migrationDir, file), version); err != nil {
			return fmt.Errorf("执行迁移 %s 失败: %w", file, err)
		}
	}

	return nil
}

// DownMigrate 执行数据库向下迁移，回滚所有已执行的迁移脚本
func DownMigrate() error {
	if DBPool == nil {
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

	// 获取已执行的迁移版本（按版本号倒序）
	executedVersions, err := getExecutedVersionsOrdered()
	if err != nil {
		return fmt.Errorf("获取已执行版本失败: %w", err)
	}

	// 执行回滚
	for _, version := range executedVersions {
		downFile := filepath.Join(migrationDir, version+".down.sql")
		if _, err := os.Stat(downFile); os.IsNotExist(err) {
			fmt.Printf("跳过回滚 %s: 下迁文件不存在\n", version)
			continue
		}

		fmt.Printf("执行回滚: %s\n", version)
		if err := executeMigrationFile(downFile, version); err != nil {
			return fmt.Errorf("回滚 %s 失败: %w", version, err)
		}

		// 删除迁移记录
		if err := deleteMigrationRecord(version); err != nil {
			return fmt.Errorf("删除迁移记录失败: %w", err)
		}
	}

	return nil
}

// createMigrationTable 创建迁移记录表
func createMigrationTable() error {
	query := `CREATE TABLE IF NOT EXISTS schema_migrations (
		version VARCHAR(255) PRIMARY KEY,
		applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`
	_, err := DBPool.Exec(query)
	return err
}

// getExecutedVersions 获取已执行的迁移版本
func getExecutedVersions() (map[string]bool, error) {
	versions := make(map[string]bool)
	rows, err := DBPool.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return versions, err
	}
	defer rows.Close()

	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return versions, err
		}
		versions[version] = true
	}
	return versions, nil
}

// getExecutedVersionsOrdered 获取已执行的迁移版本（按版本号倒序）
func getExecutedVersionsOrdered() ([]string, error) {
	var versions []string
	rows, err := DBPool.Query("SELECT version FROM schema_migrations ORDER BY version DESC")
	if err != nil {
		return versions, err
	}
	defer rows.Close()

	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return versions, err
		}
		versions = append(versions, version)
	}
	return versions, nil
}

// deleteMigrationRecord 删除迁移记录
func deleteMigrationRecord(version string) error {
	_, err := DBPool.Exec("DELETE FROM schema_migrations WHERE version = ?", version)
	return err
}

// getMigrationDir 获取迁移脚本目录的绝对路径
// 优先使用 MIGRATION_DIR 环境变量，否则回退到源码中的默认路径
func getMigrationDir() (string, error) {
	// Docker 部署时通过环境变量指定迁移脚本路径
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

	// 本地开发时使用源码中的默认路径
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
	// 从 "0001_init.up.sql" 提取 "0001_init"
	parts := strings.Split(filename, ".")
	if len(parts) > 0 {
		return parts[0]
	}
	return filename
}

// executeMigrationFile 执行单个迁移文件
func executeMigrationFile(filePath string, version string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("读取文件失败: %w", err)
	}

	sqlContent := strings.TrimSpace(string(content))
	if sqlContent == "" {
		return nil
	}

	// 开始事务
	tx, err := DBPool.Begin()
	if err != nil {
		return fmt.Errorf("开始事务失败: %w", err)
	}

	// 按分号分割并执行每条SQL
	queries := splitSQL(sqlContent)
	for _, query := range queries {
		query = strings.TrimSpace(query)
		if query == "" {
			continue
		}
		if _, err := tx.Exec(query); err != nil {
			tx.Rollback()
			return fmt.Errorf("执行SQL失败: %w\nSQL: %s", err, query)
		}
	}

	// 记录迁移版本（只对 up 文件）
	if strings.HasSuffix(filePath, ".up.sql") {
		if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES (?)", version); err != nil {
			tx.Rollback()
			return fmt.Errorf("记录迁移版本失败: %w", err)
		}
	}

	return tx.Commit()
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

		// 处理块注释 /* ... */
		if inBlockComment {
			if ch == '*' && i+1 < len(sql) && sql[i+1] == '/' {
				inBlockComment = false
				i++ // 跳过 '/'
			}
			continue
		}

		// 处理行注释 -- ...
		if inLineComment {
			if ch == '\n' {
				inLineComment = false
			} else {
				continue
			}
		}

		// 检测注释开始
		if !inSingleQuote && !inDoubleQuote {
			if ch == '-' && i+1 < len(sql) && sql[i+1] == '-' {
				inLineComment = true
				i++ // 跳过第二个 '-'
				continue
			}
			if ch == '/' && i+1 < len(sql) && sql[i+1] == '*' {
				inBlockComment = true
				i++ // 跳过 '*'
				continue
			}
		}

		// 处理引号
		if ch == '\'' && !inDoubleQuote {
			if i+1 < len(sql) && sql[i+1] == '\'' {
				// 转义的单引号
				current.WriteByte(ch)
				current.WriteByte(sql[i+1])
				i++
				continue
			}
			inSingleQuote = !inSingleQuote
		} else if ch == '"' && !inSingleQuote {
			inDoubleQuote = !inDoubleQuote
		}

		// 处理分号
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

	// 处理最后一条没有分号的SQL
	q := strings.TrimSpace(current.String())
	if q != "" {
		queries = append(queries, q)
	}

	return queries
}
