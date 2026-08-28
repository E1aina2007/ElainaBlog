package site

import (
	"ElainaBlog/internal/config"
	cache "ElainaBlog/internal/middleware/redis"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Service struct {
	repo *Repository
	rdb  *goredis.Client
	db   *gorm.DB
}

func NewService(repo *Repository, redis *goredis.Client, gormDB *gorm.DB) *Service {
	return &Service{repo: repo, rdb: redis, db: gormDB}
}

var (
	ErrDBNotInitialized = errors.New("数据库未初始化")
)

// GetDashboardStats 获取仪表盘统计数据
func (s *Service) GetDashboardStats(ctx context.Context) (*DashboardStats, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}
	stats, err := s.repo.GetDashboardStats(ctx)
	if err != nil {
		return nil, err
	}

	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	if pv, err := s.rdb.Get(ctx, "pv:"+today).Int64(); err == nil {
		stats.TodayPV = pv
	}
	if pv, err := s.rdb.Get(ctx, "pv:"+yesterday).Int64(); err == nil {
		stats.YesterdayPV = pv
	}
	if uv, err := s.rdb.SCard(ctx, "uv:"+today).Result(); err == nil {
		stats.TodayUV = uv
	}
	if uv, err := s.rdb.SCard(ctx, "uv:"+yesterday).Result(); err == nil {
		stats.YesterdayUV = uv
	}

	return stats, nil
}

// RecordVisit 记录一次页面访问（PV+UV）
func (s *Service) RecordVisit(ctx context.Context, clientIP string) {
	today := time.Now().Format("2006-01-02")

	pvKey := "pv:" + today
	s.rdb.Incr(ctx, pvKey)
	s.rdb.Expire(ctx, pvKey, 48*time.Hour)

	uvKey := "uv:" + today
	s.rdb.SAdd(ctx, uvKey, clientIP)
	s.rdb.Expire(ctx, uvKey, 48*time.Hour)
}

// GetAuthorStats 获取作者页统计数据
func (s *Service) GetAuthorStats(ctx context.Context) (*AuthorStats, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}
	return s.repo.GetAuthorStats(ctx)
}

// ExportDatabaseBackup 通过 GORM 导出数据库备份
func (s *Service) ExportDatabaseBackup(ctx context.Context) ([]byte, error) {
	if s == nil || s.db == nil {
		return nil, ErrDBNotInitialized
	}

	var out strings.Builder
	dbName := config.GlobalConfig.Db.DBName

	out.WriteString("-- ElainaBlog Database Backup\n")
	out.WriteString(fmt.Sprintf("-- Database: %s\n", dbName))
	out.WriteString(fmt.Sprintf("-- Time: %s\n\n", time.Now().Format("2006-01-02 15:04:05")))
	out.WriteString("SET FOREIGN_KEY_CHECKS=0;\n\n")

	// 获取所有表名
	rows, err := s.db.WithContext(ctx).Raw("SHOW TABLES").Rows()
	if err != nil {
		return nil, fmt.Errorf("获取表列表失败: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		tables = append(tables, name)
	}

	for _, table := range tables {
		// 获取建表 DDL
		var tableName, ddl string
		if err := s.db.WithContext(ctx).Raw(fmt.Sprintf("SHOW CREATE TABLE `%s`", table)).Row().Scan(&tableName, &ddl); err != nil {
			return nil, fmt.Errorf("获取表 %s 结构失败: %w", table, err)
		}
		out.WriteString(fmt.Sprintf("-- Table: %s\n", table))
		out.WriteString(fmt.Sprintf("DROP TABLE IF EXISTS `%s`;\n", table))
		out.WriteString(ddl + ";\n\n")

		// 读取表数据
		dataRows, err := s.db.WithContext(ctx).Raw(fmt.Sprintf("SELECT * FROM `%s`", table)).Rows()
		if err != nil {
			return nil, fmt.Errorf("读取表 %s 数据失败: %w", table, err)
		}

		cols, _ := dataRows.Columns()
		vals := make([]any, len(cols))
		ptrs := make([][]byte, len(cols))
		for i := range vals {
			vals[i] = &ptrs[i]
		}

		hasData := false
		for dataRows.Next() {
			if err := dataRows.Scan(vals...); err != nil {
				dataRows.Close()
				return nil, fmt.Errorf("扫描表 %s 行数据失败: %w", table, err)
			}
			if !hasData {
				out.WriteString(fmt.Sprintf("INSERT INTO `%s` VALUES\n", table))
				hasData = true
			} else {
				out.WriteString(",\n")
			}
			out.WriteString("(")
			for i, ptr := range ptrs {
				if i > 0 {
					out.WriteString(",")
				}
				if ptr == nil {
					out.WriteString("NULL")
				} else {
					out.WriteString(fmt.Sprintf("'%s'", strings.ReplaceAll(string(ptr), "'", "\\'")))
				}
			}
			out.WriteString(")")
		}
		dataRows.Close()

		if hasData {
			out.WriteString(";\n")
		}
		out.WriteString("\n")
	}

	out.WriteString("SET FOREIGN_KEY_CHECKS=1;\n")
	return []byte(out.String()), nil
}

// GetBannedIPs 获取被封禁的IP列表
func (s *Service) GetBannedIPs(ctx context.Context) []string {
	ips, err := cache.GetBannedIPs(s.rdb)
	if err != nil {
		return []string{}
	}
	return ips
}

// UnbanIP 解封IP
func (s *Service) UnbanIP(ctx context.Context, ip string) error {
	return cache.UnbanIP(s.rdb, ip)
}
