package article

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"ElainaBlog/config/db"
	"ElainaBlog/pkg/rdb"
)

// 返回给前端
type ArticleVO struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	Username     string    `json:"author_name"`   // 关联查询作者名
	Avatar       string    `json:"author_avatar"` // 关联查询作者头像
	IsAdmin      bool      `json:"author_is_admin"` // 关联查询作者是否管理员
	CategoryID   *int64    `json:"category_id"`   // nil 表示未分类
	CategoryName string    `json:"category_name"` // 关联查询分类名
	Title        string    `json:"title"`
	Summary      string    `json:"summary"`
	Content      string    `json:"content"`
	IsTop        bool      `json:"is_top"`
	IsDraft      bool      `json:"is_draft"`
	ViewCount    int       `json:"view_count"`
	UVCount      int64     `json:"uv_count"`
	CommentCount int       `json:"comment_count"`
	CreatedAt    time.Time `json:"created_at"`
}

// MySQLRepository 实现 article.Repository 接口，使用 MySQL 存储。
type MySQLRepository struct {
	db  db.DBTX
	rdb rdb.RedisClient // 可选，用于浏览量缓冲
}

// NewRepository 创建文章仓储实例。
// redis 参数可选，传 nil 时浏览量直接写 MySQL。
func NewRepository(db db.DBTX, redis rdb.RedisClient) *MySQLRepository {
	return &MySQLRepository{db: db, rdb: redis}
}

func (r *MySQLRepository) GetArticleByID(id int64) (*ArticleVO, error) {
	return r.getArticleByID(id, true)
}

func (r *MySQLRepository) GetArticleByIDIncludeDraft(id int64) (*ArticleVO, error) {
	return r.getArticleByID(id, false)
}

func (r *MySQLRepository) getArticleByID(id int64, filterDraft bool) (*ArticleVO, error) {
	var vo ArticleVO
	var categoryID sql.NullInt64
	var categoryName string
	query := `
		SELECT a.id, a.user_id, u.username, COALESCE(u.avatar,''), u.is_admin, a.category_id, COALESCE(c.name,''),
		       a.title, a.summary, a.content, a.is_top, a.is_draft, a.view_count, a.created_at
		FROM article a
		LEFT JOIN ` + "`user`" + ` u ON a.user_id = u.id
		LEFT JOIN category c ON a.category_id = c.id AND c.is_deleted = 0
		WHERE a.id = ? AND a.is_deleted = 0`
	if filterDraft {
		query += " AND a.is_draft = 0"
	}
	err := r.db.QueryRow(query, id).Scan(
		&vo.ID, &vo.UserID, &vo.Username, &vo.Avatar, &vo.IsAdmin, &categoryID, &categoryName,
		&vo.Title, &vo.Summary, &vo.Content, &vo.IsTop, &vo.IsDraft, &vo.ViewCount, &vo.CreatedAt)
	if err != nil {
		return nil, err
	}
	// 叠加 Redis 中的浏览量增量
	vo.ViewCount += r.GetViewCountDelta(vo.ID)
	if categoryID.Valid {
		vo.CategoryID = &categoryID.Int64
		vo.CategoryName = categoryName
	}
	return &vo, nil
}

// buildOrderBy 根据 sortBy 参数构建 ORDER BY 子句
func buildOrderBy(sortBy string) string {
	switch sortBy {
	case "popular":
		return "a.is_top DESC, a.view_count DESC"
	default:
		return "a.is_top DESC, a.created_at DESC"
	}
}

// GetArticleList 公开文章列表，过滤草稿，支持分页和分类筛选
func (r *MySQLRepository) GetArticleList(categoryID *int64, sortBy string, page, pageSize int) ([]*ArticleVO, int, error) {
	// 构建基础查询条件
	whereClause := "WHERE a.is_deleted = 0 AND a.is_draft = 0"
	args := []interface{}{}

	if categoryID != nil && *categoryID > 0 {
		whereClause += " AND a.category_id = ?"
		args = append(args, *categoryID)
	}

	// 先查询总数
	var total int
	countQuery := "SELECT COUNT(*) FROM article a " + whereClause
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 分页查询数据
	query := `
		SELECT a.id, a.user_id, u.username, COALESCE(u.avatar,''), u.is_admin, a.category_id, COALESCE(c.name,''),
		       a.title, a.summary, a.content, a.is_top, a.is_draft, a.view_count,
		       (SELECT COUNT(*) FROM comment ct WHERE ct.article_id = a.id AND ct.is_deleted = 0) AS comment_count,
		       a.created_at
		FROM article a
		LEFT JOIN ` + "`user`" + ` u ON a.user_id = u.id
		LEFT JOIN category c ON a.category_id = c.id AND c.is_deleted = 0
		` + whereClause + `
		ORDER BY ` + buildOrderBy(sortBy) + `
		LIMIT ? OFFSET ?`

	offset := (page - 1) * pageSize
	args = append(args, pageSize, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	articles := make([]*ArticleVO, 0)
	for rows.Next() {
		var vo ArticleVO
		var catID sql.NullInt64
		var categoryName string
		err := rows.Scan(&vo.ID, &vo.UserID, &vo.Username, &vo.Avatar, &vo.IsAdmin, &catID, &categoryName,
			&vo.Title, &vo.Summary, &vo.Content, &vo.IsTop, &vo.IsDraft, &vo.ViewCount, &vo.CommentCount, &vo.CreatedAt)
		if err != nil {
			return nil, 0, err
		}
		if catID.Valid {
			vo.CategoryID = &catID.Int64
			vo.CategoryName = categoryName
		}
		articles = append(articles, &vo)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return articles, total, nil
}

// GetAdminArticleList 管理员文章列表，包含草稿，支持分页和分类筛选
func (r *MySQLRepository) GetAdminArticleList(categoryID *int64, sortBy string, page, pageSize int) ([]*ArticleVO, int, error) {
	whereClause := "WHERE a.is_deleted = 0"
	args := []interface{}{}

	if categoryID != nil && *categoryID > 0 {
		whereClause += " AND a.category_id = ?"
		args = append(args, *categoryID)
	}

	var total int
	countQuery := "SELECT COUNT(*) FROM article a " + whereClause
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT a.id, a.user_id, u.username, COALESCE(u.avatar,''), u.is_admin, a.category_id, COALESCE(c.name,''),
		       a.title, a.summary, a.content, a.is_top, a.is_draft, a.view_count,
		       (SELECT COUNT(*) FROM comment ct WHERE ct.article_id = a.id AND ct.is_deleted = 0) AS comment_count,
		       a.created_at
		FROM article a
		LEFT JOIN ` + "`user`" + ` u ON a.user_id = u.id
		LEFT JOIN category c ON a.category_id = c.id AND c.is_deleted = 0
		` + whereClause + `
		ORDER BY ` + buildOrderBy(sortBy) + `
		LIMIT ? OFFSET ?`

	offset := (page - 1) * pageSize
	args = append(args, pageSize, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	articles := make([]*ArticleVO, 0)
	for rows.Next() {
		var vo ArticleVO
		var catID sql.NullInt64
		var categoryName string
		err := rows.Scan(&vo.ID, &vo.UserID, &vo.Username, &vo.Avatar, &vo.IsAdmin, &catID, &categoryName,
			&vo.Title, &vo.Summary, &vo.Content, &vo.IsTop, &vo.IsDraft, &vo.ViewCount, &vo.CommentCount, &vo.CreatedAt)
		if err != nil {
			return nil, 0, err
		}
		if catID.Valid {
			vo.CategoryID = &catID.Int64
			vo.CategoryName = categoryName
		}
		articles = append(articles, &vo)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	// 填充每篇文章的 UV 数
	for _, vo := range articles {
		if uv, err := r.GetArticleUV(vo.ID); err == nil {
			vo.UVCount = uv
		}
	}

	return articles, total, nil
}

// GetUserArticleList 用户自己的文章列表（含草稿），支持分页和分类筛选
func (r *MySQLRepository) GetUserArticleList(userID int64, categoryID *int64, sortBy string, page, pageSize int) ([]*ArticleVO, int, error) {
	whereClause := "WHERE a.is_deleted = 0 AND a.user_id = ?"
	args := []interface{}{userID}

	if categoryID != nil && *categoryID > 0 {
		whereClause += " AND a.category_id = ?"
		args = append(args, *categoryID)
	}

	var total int
	countQuery := "SELECT COUNT(*) FROM article a " + whereClause
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT a.id, a.user_id, u.username, COALESCE(u.avatar,''), u.is_admin, a.category_id, COALESCE(c.name,''),
		       a.title, a.summary, a.content, a.is_top, a.is_draft, a.view_count,
		       (SELECT COUNT(*) FROM comment ct WHERE ct.article_id = a.id AND ct.is_deleted = 0) AS comment_count,
		       a.created_at
		FROM article a
		LEFT JOIN ` + "`user`" + ` u ON a.user_id = u.id
		LEFT JOIN category c ON a.category_id = c.id AND c.is_deleted = 0
		` + whereClause + `
		ORDER BY ` + buildOrderBy(sortBy) + `
		LIMIT ? OFFSET ?`

	offset := (page - 1) * pageSize
	args = append(args, pageSize, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	articles := make([]*ArticleVO, 0)
	for rows.Next() {
		var vo ArticleVO
		var catID sql.NullInt64
		var categoryName string
		err := rows.Scan(&vo.ID, &vo.UserID, &vo.Username, &vo.Avatar, &vo.IsAdmin, &catID, &categoryName,
			&vo.Title, &vo.Summary, &vo.Content, &vo.IsTop, &vo.IsDraft, &vo.ViewCount, &vo.CommentCount, &vo.CreatedAt)
		if err != nil {
			return nil, 0, err
		}
		if catID.Valid {
			vo.CategoryID = &catID.Int64
			vo.CategoryName = categoryName
		}
		articles = append(articles, &vo)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return articles, total, nil
}

// IncrementViewCount 增加文章浏览量（Redis 缓冲，定时同步到 MySQL）
func (r *MySQLRepository) IncrementViewCount(id int64) error {
	// Redis 可用时写入 Redis 缓冲
	if r.rdb != nil {
		ctx := context.Background()
		key := fmt.Sprintf("article:view_count:%d", id)
		err := r.rdb.Incr(ctx, key).Err()
		if err == nil {
			return nil
		}
		// Redis 故障时 fallback 到直接写 MySQL
	}

	_, err := r.db.Exec("UPDATE article SET view_count = view_count + 1 WHERE id = ? AND is_deleted = 0", id)
	return err
}

// IncrementViewCountUnique 带 IP 去重的浏览量递增
// 使用 Redis Set 记录每篇文章的访问 IP，同一 IP 永久只计一次
func (r *MySQLRepository) IncrementViewCountUnique(id int64, clientIP string) error {
	if r.rdb != nil {
		ctx := context.Background()
		uvKey := fmt.Sprintf("article:uv:%d", id)

		// 检查该 IP 是否已访问过此文章
		exists, err := r.rdb.SIsMember(ctx, uvKey, clientIP).Result()
		if err == nil && exists {
			return nil // 已访问过，跳过计数
		}

		// 首次访问，加入 Set（永久保留，不设 TTL）
		r.rdb.SAdd(ctx, uvKey, clientIP)
	}

	// 走原有的 INCR 逻辑
	return r.IncrementViewCount(id)
}

// GetArticleUV 获取文章的独立访客数（UV）
func (r *MySQLRepository) GetArticleUV(id int64) (int64, error) {
	if r.rdb == nil {
		return 0, nil
	}
	ctx := context.Background()
	uvKey := fmt.Sprintf("article:uv:%d", id)
	return r.rdb.SCard(ctx, uvKey).Result()
}

// GetViewCountDelta 从 Redis 读取文章浏览量增量（不删除，由 FlushViewCounts 定时清理）
func (r *MySQLRepository) GetViewCountDelta(id int64) int {
	if r.rdb == nil {
		return 0
	}
	ctx := context.Background()
	key := fmt.Sprintf("article:view_count:%d", id)
	val, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		return 0
	}
	var delta int
	fmt.Sscanf(val, "%d", &delta)
	return delta
}

// FlushViewCounts 将 Redis 中累积的浏览量批量同步到 MySQL
func (r *MySQLRepository) FlushViewCounts() (int, error) {
	if r.rdb == nil {
		return 0, nil
	}

	ctx := context.Background()

	var cursor uint64
	flushed := 0
	for {
		keys, nextCursor, err := r.rdb.Scan(ctx, cursor, "article:view_count:*", 100).Result()
		if err != nil {
			return flushed, err
		}

		for _, key := range keys {
			val, err := r.rdb.GetDel(ctx, key).Result()
			if err != nil {
				continue
			}

			var delta int
			fmt.Sscanf(val, "%d", &delta)
			if delta <= 0 {
				continue
			}

			var articleID int64
			fmt.Sscanf(key, "article:view_count:%d", &articleID)
			if articleID <= 0 {
				continue
			}

			_, err = r.db.Exec("UPDATE article SET view_count = view_count + ? WHERE id = ? AND is_deleted = 0", delta, articleID)
			if err != nil {
				continue
			}
			flushed++
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return flushed, nil
}

func (r *MySQLRepository) CreateArticle(userID int64, categoryID *int64, title, summary, content string, isTop, isDraft bool) (int64, error) {
	result, err := r.db.Exec("INSERT INTO article (user_id, category_id, title, summary, content, is_top, is_draft) VALUES (?, ?, ?, ?, ?, ?, ?)",
		userID, categoryID, title, summary, content, isTop, isDraft)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *MySQLRepository) UpdateArticle(id int64, categoryID *int64, title, summary, content string, isTop, isDraft bool) error {
	_, err := r.db.Exec("UPDATE article SET category_id = ?, title = ?, summary = ?, content = ?, is_top = ?, is_draft = ? WHERE id = ? AND is_deleted = 0",
		categoryID, title, summary, content, isTop, isDraft, id)
	return err
}

func (r *MySQLRepository) DeleteArticle(id int64) error {
	_, err := r.db.Exec("UPDATE article SET is_deleted = 1 WHERE id = ? AND is_deleted = 0", id)
	return err
}

// GetArticleAuthorInfo 获取文章作者 ID 和标题（供通知使用）
func (r *MySQLRepository) GetArticleAuthorInfo(id int64) (userID int64, title string, err error) {
	err = r.db.QueryRow("SELECT user_id, title FROM article WHERE id = ? AND is_deleted = 0", id).Scan(&userID, &title)
	return
}
