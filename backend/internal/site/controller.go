package site

import (
	"ElainaBlog/config/db"
	"ElainaBlog/internal/common"
	"ElainaBlog/internal/common/model"
	"ElainaBlog/internal/user"
	"ElainaBlog/pkg/rdb"
	"bufio"
	"context"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service     *Service
	userService *user.Service
}

func NewController(service *Service, userService *user.Service) *Controller {
	return &Controller{service: service, userService: userService}
}

// GetDashboardStats 获取仪表盘统计数据（管理员）
func (ctl *Controller) GetDashboardStats(c *gin.Context) {
	if !common.RequireAdmin(c, ctl.userService.CheckIsAdmin) {
		return
	}
	stats, err := ctl.service.GetDashboardStats()
	if err != nil {
		c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
		return
	}
	c.JSON(http.StatusOK, model.ApiSuccessResponse(stats))
}

// GetAuthorStats 获取作者页统计数据（公开接口）
func (ctl *Controller) GetAuthorStats(c *gin.Context) {
	stats, err := ctl.service.GetAuthorStats()
	if err != nil {
		c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
		return
	}
	c.JSON(http.StatusOK, model.ApiSuccessResponse(stats))
}

// SystemStatus 系统状态响应结构
type SystemStatusResponse struct {
	CPUUsage     float64 `json:"cpu_usage"`
	MemoryUsage  float64 `json:"memory_usage"`
	MemoryTotal  uint64  `json:"memory_total"`
	MemoryUsed   uint64  `json:"memory_used"`
	DBStatus     string  `json:"db_status"`
	RedisStatus  string  `json:"redis_status"`
	CacheHitRate float64 `json:"cache_hit_rate"`
	Uptime       string  `json:"uptime"`
}

// cpuTimes 存储上一次 CPU 时间快照
var (
	lastCPUTimes  [4]uint64
	lastCPULock   sync.Mutex
	lastCPUUsage  float64
)

// readCPUTimes 从 /proc/stat 读取 CPU 时间（user, nice, system, idle）
func readCPUTimes() ([4]uint64, error) {
	f, err := os.Open("/proc/stat")
	if err != nil {
		return [4]uint64{}, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "cpu ") {
			fields := strings.Fields(line)
			if len(fields) < 5 {
				return [4]uint64{}, fmt.Errorf("unexpected /proc/stat format")
			}
			var times [4]uint64
			for i := 0; i < 4; i++ {
				v, err := strconv.ParseUint(fields[i+1], 10, 64)
				if err != nil {
					return [4]uint64{}, err
				}
				times[i] = v
			}
			return times, nil
		}
	}
	return [4]uint64{}, fmt.Errorf("cpu line not found in /proc/stat")
}

// getCPUUsage 计算 CPU 使用率（两次采样取差值）
func getCPUUsage() float64 {
	times, err := readCPUTimes()
	if err != nil {
		return 0
	}

	lastCPULock.Lock()
	defer lastCPULock.Unlock()

	prev := lastCPUTimes
	lastCPUTimes = times

	prevTotal := prev[0] + prev[1] + prev[2] + prev[3]
	total := times[0] + times[1] + times[2] + times[3]

	if total == prevTotal || prevTotal == 0 {
		return lastCPUUsage
	}

	prevIdle := prev[3]
	idle := times[3]

	totalDiff := total - prevTotal
	idleDiff := idle - prevIdle

	usage := float64(totalDiff-idleDiff) / float64(totalDiff) * 100
	lastCPUUsage = float64(int(usage*100+0.5)) / 100
	return lastCPUUsage
}

// readMemInfo 从 /proc/meminfo 读取系统内存信息
func readMemInfo() (totalMB, usedMB uint64, usagePercent float64) {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0, 0, 0
	}
	defer f.Close()

	var memTotal, memAvailable uint64
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "MemTotal:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				v, _ := strconv.ParseUint(fields[1], 10, 64)
				memTotal = v // kB
			}
		} else if strings.HasPrefix(line, "MemAvailable:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				v, _ := strconv.ParseUint(fields[1], 10, 64)
				memAvailable = v // kB
			}
		}
		if memTotal > 0 && memAvailable > 0 {
			break
		}
	}

	totalMB = memTotal / 1024
	usedMB = (memTotal - memAvailable) / 1024
	if memTotal > 0 {
		usagePercent = float64(int(float64(memTotal-memAvailable)/float64(memTotal)*10000+0.5)) / 100
	}
	return
}

// GetSystemStatus 获取系统运行状态（管理员）
func (ctl *Controller) GetSystemStatus(c *gin.Context) {
	if !common.RequireAdmin(c, ctl.userService.CheckIsAdmin) {
		return
	}

	// 检查数据库连接
	dbStatus := "connected"
	if err := db.DBPool.Ping(); err != nil {
		dbStatus = "error"
	}

	// 检查 Redis 连接并获取缓存命中率
	redisStatus := "connected"
	cacheHitRate := float64(-1)
	if rdb.RedisClient == nil {
		redisStatus = "not_initialized"
	} else if err := rdb.RedisClient.Ping(context.Background()).Err(); err != nil {
		redisStatus = "error"
	} else {
		cacheHitRate = getRedisCacheHitRate()
	}

	// CPU 使用率（从 /proc/stat 读取）
	cpuUsage := getCPUUsage()

	// 内存使用率（从 /proc/meminfo 读取）
	memTotal, memUsed, memUsage := readMemInfo()
	// 回退到 Go 运行时内存
	if memTotal == 0 {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		memTotal = m.Sys / 1024 / 1024
		memUsed = m.Alloc / 1024 / 1024
		if m.Sys > 0 {
			memUsage = float64(int(float64(m.Alloc)/float64(m.Sys)*10000+0.5)) / 100
		}
	}

	status := SystemStatusResponse{
		CPUUsage:     cpuUsage,
		MemoryUsage:  memUsage,
		MemoryTotal:  memTotal,
		MemoryUsed:   memUsed,
		DBStatus:     dbStatus,
		RedisStatus:  redisStatus,
		CacheHitRate: cacheHitRate,
		Uptime:       formatUptime(time.Since(startTime)),
	}

	c.JSON(http.StatusOK, model.ApiSuccessResponse(status))
}

// ClearCache 清理缓存（管理员）
func (ctl *Controller) ClearCache(c *gin.Context) {
	if !common.RequireAdmin(c, ctl.userService.CheckIsAdmin) {
		return
	}

	// 清理 Go 内存缓存
	runtime.GC()

	c.JSON(http.StatusOK, model.ApiSuccessResponse(gin.H{"message": "缓存已清理"}))
}

// ExportBackup 导出数据库备份（管理员）
func (ctl *Controller) ExportBackup(c *gin.Context) {
	if !common.RequireAdmin(c, ctl.userService.CheckIsAdmin) {
		return
	}

	// 使用 mysqldump 命令导出数据库
	// 简单实现：查询所有表数据并返回
	backup, err := ctl.service.ExportDatabaseBackup()
	if err != nil {
		c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
		return
	}

	filename := fmt.Sprintf("backup_%s.sql", time.Now().Format("2006-01-02"))
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/sql")
	c.Data(http.StatusOK, "application/sql", backup)
}

// GetBannedIPs 获取被封禁的IP列表（管理员）
func (ctl *Controller) GetBannedIPs(c *gin.Context) {
	if !common.RequireAdmin(c, ctl.userService.CheckIsAdmin) {
		return
	}

	ips := ctl.service.GetBannedIPs()
	c.JSON(http.StatusOK, model.ApiSuccessResponse(ips))
}

// UnbanIP 解封IP（管理员）
func (ctl *Controller) UnbanIP(c *gin.Context) {
	if !common.RequireAdmin(c, ctl.userService.CheckIsAdmin) {
		return
	}

	var req struct {
		IP string `json:"ip"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.IP == "" {
		appErr := model.ErrInvalidParams.WithDetail("无效的IP地址")
		c.JSON(appErr.HTTPStatus(), model.ApiErrorResponse(appErr.Code, appErr.Message, appErr))
		return
	}

	if err := ctl.service.UnbanIP(req.IP); err != nil {
		c.JSON(model.ErrInternal.HTTPStatus(), model.ApiErrorResponse(model.ErrInternal.Code, model.ErrInternal.Message, nil))
		return
	}
	c.JSON(http.StatusOK, model.ApiSuccessResponse(gin.H{"message": "IP已解封"}))
}

// startTime 记录服务启动时间
var startTime = time.Now()

// getRedisCacheHitRate 从 Redis INFO stats 中计算缓存命中率
func getRedisCacheHitRate() float64 {
	info, err := rdb.RedisClient.Info(context.Background(), "stats").Result()
	if err != nil {
		return -1
	}

	var hits, misses float64
	for _, line := range strings.Split(info, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "keyspace_hits:") {
			hits, _ = strconv.ParseFloat(strings.TrimPrefix(line, "keyspace_hits:"), 64)
		} else if strings.HasPrefix(line, "keyspace_misses:") {
			misses, _ = strconv.ParseFloat(strings.TrimPrefix(line, "keyspace_misses:"), 64)
		}
	}

	total := hits + misses
	if total == 0 {
		return 0
	}
	return float64(int(hits/total*10000+0.5)) / 100
}

// formatUptime 格式化运行时间，只显示分及以上单位
func formatUptime(d time.Duration) string {
	totalMinutes := int(d.Minutes())
	days := totalMinutes / (24 * 60)
	hours := (totalMinutes % (24 * 60)) / 60
	minutes := totalMinutes % 60
	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	}
	return fmt.Sprintf("%dh %dm", hours, minutes)
}
