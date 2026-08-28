package site

// DashboardStats 仪表盘统计数据
type DashboardStats struct {
	ArticleCount int64 `json:"article_count"`
	CommentCount int64 `json:"comment_count"`
	UserCount    int64 `json:"user_count"`
	TodayPV      int64 `json:"today_pv"`
	TodayUV      int64 `json:"today_uv"`
	YesterdayPV  int64 `json:"yesterday_pv"`
	YesterdayUV  int64 `json:"yesterday_uv"`
}

// AuthorStats 关于作者页统计数据
type AuthorStats struct {
	ArticleCount int64 `json:"article_count"`
	CommentCount int64 `json:"comment_count"`
	TotalViews   int64 `json:"total_views"`
	DaysSince    int   `json:"days_since"`
}

// SystemStatusResponse 系统状态响应结构
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
