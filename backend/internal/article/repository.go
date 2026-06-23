package article

import (
	"context"
	"fmt"
	"time"

	"ElainaBlog/pkg/rdb"

	"gorm.io/gorm"
)

// 返回给前端
type ArticleVO struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	Username     string    `json:"author_name"`
	Avatar       string    `json:"author_avatar"`
	IsAdmin      bool      `json:"author_is_admin"`
	CategoryID   *int64    `json:"category_id"`
	CategoryName string    `json:"category_name"`
	Title        string    `json:"title"`
	Summary      string    `json:"summary"`
	Content      string    `json:"content"`
	Tags         string    `json:"tags"`
	IsTop        bool      `json:"is_top"`
	IsDraft      bool      `json:"is_draft"`
	ViewCount    int       `json:"view_count"`
	UVCount      int64     `json:"uv_count"`
	CommentCount int       `json:"comment_count"`
	CreatedAt    time.Time `json:"created_at"`
}

// MySQLRepository 实现 article.Repository 接口，使用 GORM 存储。
type MySQLRepository struct {
	db  *gorm.DB
	rdb rdb.RedisClient // 可选，用于浏览量缓冲
}

// NewRepository 创建文章仓储实例。
func NewRepository(db *gorm.DB, redis rdb.RedisClient) *MySQLRepository {
	return &MySQLRepository{db: db, rdb: redis}
}

// articleSelect 返回文章列表查询的公共 Select 子句
const articleSelect = `
	a.id, a.user_id, u.username, COALESCE(u.avatar,'') AS author_avatar, u.is_admin AS author_is_admin,
	a.category_id, COALESCE(c.name,'') AS category_name,
	a.title, a.summary, a.content, COALESCE(a.tags,'') AS tags, a.is_top, a.is_draft, a.view_count,
	(SELECT COUNT(*) FROM comment ct WHERE ct.article_id = a.id AND ct.is_deleted = 0) AS comment_count,
	a.created_at`

// articleJoins 返回文章列表查询的公共 Join 子句
func articleJoins(db *gorm.DB) *gorm.DB {
	return db.Joins("LEFT JOIN `user` u ON a.user_id = u.id").
		Joins("LEFT JOIN category c ON a.category_id = c.id AND c.is_deleted = 0")
}

func (r *MySQLRepository) GetArticleByID(id int64) (*ArticleVO, error) {
	return r.getArticleByID(id, true)
}

func (r *MySQLRepository) GetArticleByIDIncludeDraft(id int64) (*ArticleVO, error) {
	return r.getArticleByID(id, false)
}

func (r *MySQLRepository) getArticleByID(id int64, filterDraft bool) (*ArticleVO, error) {
	var vo ArticleVO
	query := r.db.Table("article a").Select(articleSelect)
	query = articleJoins(query)
	query = query.Where("a.id = ? AND a.is_deleted = 0", id)
	if filterDraft {
		query = query.Where("a.is_draft = 0")
	}
	if err := query.Scan(&vo).Error; err != nil {
		return nil, err
	}
	// 叠加 Redis 中的浏览量增量
	vo.ViewCount += r.GetViewCountDelta(vo.ID)
	return &vo, nil
}

// buildOrderBy 根据 sortBy 参数构建 ORDER BY 子句
func buildOrderBy(sortBy string) string {
	switch sortBy {
	case "popular":
		return "a.is_top DESC, a.view_count DESC"
	case "comment":
		return "a.is_top DESC, (SELECT COUNT(*) FROM comment ct WHERE ct.article_id = a.id AND ct.is_deleted = 0) DESC"
	default:
		return "a.is_top DESC, a.created_at DESC"
	}
}

// listQuery 构建文章列表查询的公共部分
func (r *MySQLRepository) listQuery(categoryID *int64, includeDraft bool, userID *int64) *gorm.DB {
	query := r.db.Table("article a").Select(articleSelect)
	query = articleJoins(query)

	query = query.Where("a.is_deleted = 0")
	if !includeDraft {
		query = query.Where("a.is_draft = 0")
	}
	if userID != nil {
		query = query.Where("a.user_id = ?", *userID)
	}
	if categoryID != nil && *categoryID > 0 {
		query = query.Where("a.category_id = ?", *categoryID)
	}
	return query
}

// countQuery 构建文章计数查询
func (r *MySQLRepository) countQuery(categoryID *int64, includeDraft bool, userID *int64) *gorm.DB {
	query := r.db.Table("article a")
	query = query.Where("a.is_deleted = 0")
	if !includeDraft {
		query = query.Where("a.is_draft = 0")
	}
	if userID != nil {
		query = query.Where("a.user_id = ?", *userID)
	}
	if categoryID != nil && *categoryID > 0 {
		query = query.Where("a.category_id = ?", *categoryID)
	}
	return query
}

// GetArticleList 公开文章列表，过滤草稿，支持分页和分类筛选
func (r *MySQLRepository) GetArticleList(categoryID *int64, sortBy string, page, pageSize int) ([]*ArticleVO, int, error) {
	var total int64
	if err := r.countQuery(categoryID, false, nil).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var articles []*ArticleVO
	query := r.listQuery(categoryID, false, nil)
	err := query.
		Order(buildOrderBy(sortBy)).
		Limit(pageSize).Offset((page - 1) * pageSize).
		Scan(&articles).Error
	return articles, int(total), err
}

// GetAdminArticleList 管理员文章列表，包含草稿，支持分页和分类筛选
func (r *MySQLRepository) GetAdminArticleList(categoryID *int64, sortBy string, page, pageSize int) ([]*ArticleVO, int, error) {
	var total int64
	if err := r.countQuery(categoryID, true, nil).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var articles []*ArticleVO
	query := r.listQuery(categoryID, true, nil)
	if err := query.
		Order(buildOrderBy(sortBy)).
		Limit(pageSize).Offset((page - 1) * pageSize).
		Scan(&articles).Error; err != nil {
		return nil, 0, err
	}

	// 填充每篇文章的 UV 数
	for _, vo := range articles {
		if uv, err := r.GetArticleUV(vo.ID); err == nil {
			vo.UVCount = uv
		}
	}

	return articles, int(total), nil
}

// GetUserArticleList 用户自己的文章列表（含草稿），支持分页和分类筛选
func (r *MySQLRepository) GetUserArticleList(userID int64, categoryID *int64, sortBy string, page, pageSize int) ([]*ArticleVO, int, error) {
	var total int64
	if err := r.countQuery(categoryID, true, &userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var articles []*ArticleVO
	query := r.listQuery(categoryID, true, &userID)
	err := query.
		Order(buildOrderBy(sortBy)).
		Limit(pageSize).Offset((page - 1) * pageSize).
		Scan(&articles).Error
	return articles, int(total), err
}

// SearchArticleList 全文搜索文章列表，按相关性排序
func (r *MySQLRepository) SearchArticleList(keyword string, page, pageSize int) ([]*ArticleVO, int, error) {
	var total int64
	if err := r.db.Table("article a").
		Where("a.is_deleted = 0 AND a.is_draft = 0 AND MATCH(a.title, a.summary) AGAINST(? IN BOOLEAN MODE)", keyword).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var articles []*ArticleVO
	query := r.db.Table("article a").Select(articleSelect)
	query = articleJoins(query)
	err := query.
		Where("a.is_deleted = 0 AND a.is_draft = 0 AND MATCH(a.title, a.summary) AGAINST(? IN BOOLEAN MODE)", keyword).
		Order(fmt.Sprintf("MATCH(a.title, a.summary) AGAINST('%s' IN BOOLEAN MODE) DESC, a.created_at DESC", keyword)).
		Limit(pageSize).Offset((page - 1) * pageSize).
		Scan(&articles).Error
	return articles, int(total), err
}

// IncrementViewCount 增加文章浏览量（Redis 缓冲，定时同步到 MySQL）
func (r *MySQLRepository) IncrementViewCount(id int64) error {
	if r.rdb != nil {
		ctx := context.Background()
		key := fmt.Sprintf("article:view_count:%d", id)
		err := r.rdb.Incr(ctx, key).Err()
		if err == nil {
			return nil
		}
	}
	return r.db.Exec("UPDATE article SET view_count = view_count + 1 WHERE id = ? AND is_deleted = 0", id).Error
}

// IncrementViewCountUnique 带 IP 去重的浏览量递增
func (r *MySQLRepository) IncrementViewCountUnique(id int64, clientIP string) error {
	if r.rdb != nil {
		ctx := context.Background()
		uvKey := fmt.Sprintf("article:uv:%d", id)
		exists, err := r.rdb.SIsMember(ctx, uvKey, clientIP).Result()
		if err == nil && exists {
			return nil
		}
		r.rdb.SAdd(ctx, uvKey, clientIP)
	}
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

// GetViewCountDelta 从 Redis 读取文章浏览量增量
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

			if err := r.db.Exec("UPDATE article SET view_count = view_count + ? WHERE id = ? AND is_deleted = 0", delta, articleID).Error; err != nil {
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

func (r *MySQLRepository) CreateArticle(userID int64, categoryID *int64, title, summary, content, tags string, isTop, isDraft bool) (int64, error) {
	result := r.db.Exec("INSERT INTO article (user_id, category_id, title, summary, content, tags, is_top, is_draft) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		userID, categoryID, title, summary, content, tags, isTop, isDraft)
	if result.Error != nil {
		return 0, result.Error
	}
	var id int64
	if err := r.db.Raw("SELECT LAST_INSERT_ID()").Scan(&id).Error; err != nil {
		return 0, err
	}
	return id, nil
}

func (r *MySQLRepository) UpdateArticle(id int64, categoryID *int64, title, summary, content, tags string, isTop, isDraft bool) error {
	return r.db.Exec("UPDATE article SET category_id = ?, title = ?, summary = ?, content = ?, tags = ?, is_top = ?, is_draft = ? WHERE id = ? AND is_deleted = 0",
		categoryID, title, summary, content, tags, isTop, isDraft, id).Error
}

// ToggleArticleTop 单独切换文章置顶状态
func (r *MySQLRepository) ToggleArticleTop(id int64, isTop bool) error {
	return r.db.Exec("UPDATE article SET is_top = ? WHERE id = ? AND is_deleted = 0", isTop, id).Error
}

func (r *MySQLRepository) DeleteArticle(id int64) error {
	return r.db.Exec("UPDATE article SET is_deleted = 1 WHERE id = ? AND is_deleted = 0", id).Error
}

// GetArticleAuthorInfo 获取文章作者 ID 和标题（供通知使用）
func (r *MySQLRepository) GetArticleAuthorInfo(id int64) (userID int64, title string, err error) {
	err = r.db.Table("article").
		Select("user_id", "title").
		Where("id = ? AND is_deleted = 0", id).
		Row().Scan(&userID, &title)
	return
}

// GetAllActiveContents 获取所有未删除文章的 content（供图片清理服务使用）
func (r *MySQLRepository) GetAllActiveContents() ([]string, error) {
	var contents []string
	err := r.db.Table("article").
		Where("is_deleted = 0").
		Pluck("content", &contents).Error
	return contents, err
}
