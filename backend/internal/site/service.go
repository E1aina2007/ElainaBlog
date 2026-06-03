package site

import (
	"ElainaBlog/config"
	"ElainaBlog/pkg/rdb"
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"time"
)

type Service struct {
	repo Repository
	rdb  rdb.RedisClient // 可选，用于 IP 封禁等操作
}

func NewService(repo Repository, redis rdb.RedisClient) *Service {
	return &Service{repo: repo, rdb: redis}
}

var (
	ErrDBNotInitialized = errors.New("数据库未初始化")
)

// GetDashboardStats 获取仪表盘统计数据
func (s *Service) GetDashboardStats() (*DashboardStats, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}
	stats, err := s.repo.GetDashboardStats()
	if err != nil {
		return nil, err
	}

	// 从 Redis 读取 PV/UV 统计
	ctx := context.Background()
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
func (s *Service) RecordVisit(clientIP string) {
	ctx := context.Background()
	today := time.Now().Format("2006-01-02")

	// PV: 累加计数，TTL 48h
	pvKey := "pv:" + today
	s.rdb.Incr(ctx, pvKey)
	s.rdb.Expire(ctx, pvKey, 48*time.Hour)

	// UV: 记录 IP 到集合，TTL 48h
	uvKey := "uv:" + today
	s.rdb.SAdd(ctx, uvKey, clientIP)
	s.rdb.Expire(ctx, uvKey, 48*time.Hour)
}

// GetAuthorStats 获取作者页统计数据
func (s *Service) GetAuthorStats() (*AuthorStats, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}
	return s.repo.GetAuthorStats()
}

// ExportDatabaseBackup 通过 mysqldump 导出数据库备份
func (s *Service) ExportDatabaseBackup() ([]byte, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}

	dbCfg := config.GlobalConfig.Db

	cmd := exec.Command("mysqldump",
		"-h", dbCfg.Host,
		"-P", fmt.Sprintf("%d", dbCfg.Port),
		"-u", dbCfg.Username,
		fmt.Sprintf("-p%s", dbCfg.Password),
		"--single-transaction",
		"--routines",
		"--triggers",
		"--set-gtid-purged=OFF",
		dbCfg.DBName,
	)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("mysqldump 执行失败: %s", stderr.String())
	}

	return out.Bytes(), nil
}

// GetBannedIPs 获取被封禁的IP列表
func (s *Service) GetBannedIPs() []string {
	ips, err := rdb.GetBannedIPs(s.rdb)
	if err != nil {
		return []string{}
	}
	return ips
}

// UnbanIP 解封IP
func (s *Service) UnbanIP(ip string) error {
	return rdb.UnbanIP(s.rdb, ip)
}

