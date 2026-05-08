package article

import (
	"database/sql"
	"time"
)

// 返回给前端
type ArticleVO struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	Username     string    `json:"author_name"`   // 关联查询作者名
	Avatar       string    `json:"author_avatar"` // 关联查询作者头像
	CategoryID   *int64    `json:"category_id"`   // nil 表示未分类
	CategoryName string    `json:"category_name"` // 关联查询分类名
	Title        string    `json:"title"`
	Summary      string    `json:"summary"`
	Content      string    `json:"content"`
	Cover        string    `json:"cover"`
	IsTop        bool      `json:"is_top"`
	IsDraft      bool      `json:"is_draft"`
	ViewCount    int       `json:"view_count"`
	CommentCount int       `json:"comment_count"`
	CreatedAt    time.Time `json:"created_at"`
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetArticleByID(id int64) (*ArticleVO, error) {
	var vo ArticleVO
	var categoryID sql.NullInt64
	var categoryName string
	err := r.db.QueryRow(`
		SELECT a.id, a.user_id, u.username, COALESCE(u.avatar,''), a.category_id, COALESCE(c.name,''),
		       a.title, a.summary, a.content, a.cover, a.is_top, a.is_draft, a.view_count, a.created_at
		FROM article a
		LEFT JOIN `+"`user`"+` u ON a.user_id = u.id
		LEFT JOIN category c ON a.category_id = c.id AND c.is_deleted = 0
		WHERE a.id = ? AND a.is_deleted = 0`, id).Scan(
		&vo.ID, &vo.UserID, &vo.Username, &vo.Avatar, &categoryID, &categoryName,
		&vo.Title, &vo.Summary, &vo.Content, &vo.Cover, &vo.IsTop, &vo.IsDraft, &vo.ViewCount, &vo.CreatedAt)
	if err != nil {
		return nil, err
	}
	if categoryID.Valid {
		vo.CategoryID = &categoryID.Int64
		vo.CategoryName = categoryName
	}
	return &vo, nil
}

// GetArticleList 公开文章列表，过滤草稿，支持分页和分类筛选
func (r *Repository) GetArticleList(categoryID *int64, page, pageSize int) ([]*ArticleVO, int, error) {
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
		SELECT a.id, a.user_id, u.username, COALESCE(u.avatar,''), a.category_id, COALESCE(c.name,''),
		       a.title, a.summary, a.content, a.cover, a.is_top, a.is_draft, a.view_count,
		       (SELECT COUNT(*) FROM comment ct WHERE ct.article_id = a.id AND ct.is_deleted = 0) AS comment_count,
		       a.created_at
		FROM article a
		LEFT JOIN ` + "`user`" + ` u ON a.user_id = u.id
		LEFT JOIN category c ON a.category_id = c.id AND c.is_deleted = 0
		` + whereClause + `
		ORDER BY a.is_top DESC, a.created_at DESC
		LIMIT ? OFFSET ?`

	offset := (page - 1) * pageSize
	args = append(args, pageSize, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var articles []*ArticleVO
	for rows.Next() {
		var vo ArticleVO
		var catID sql.NullInt64
		var categoryName string
		err := rows.Scan(&vo.ID, &vo.UserID, &vo.Username, &vo.Avatar, &catID, &categoryName,
			&vo.Title, &vo.Summary, &vo.Content, &vo.Cover, &vo.IsTop, &vo.IsDraft, &vo.ViewCount, &vo.CommentCount, &vo.CreatedAt)
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

// IncrementViewCount 增加文章浏览量
func (r *Repository) IncrementViewCount(id int64) error {
	_, err := r.db.Exec("UPDATE article SET view_count = view_count + 1 WHERE id = ? AND is_deleted = 0", id)
	return err
}

func (r *Repository) CreateArticle(userID int64, categoryID *int64, title, summary, content, cover string, isTop, isDraft bool) (int64, error) {
	result, err := r.db.Exec("INSERT INTO article (user_id, category_id, title, summary, content, cover, is_top, is_draft) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		userID, categoryID, title, summary, content, cover, isTop, isDraft)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *Repository) UpdateArticle(id int64, categoryID *int64, title, summary, content, cover string, isTop, isDraft bool) error {
	_, err := r.db.Exec("UPDATE article SET category_id = ?, title = ?, summary = ?, content = ?, cover = ?, is_top = ?, is_draft = ? WHERE id = ? AND is_deleted = 0",
		categoryID, title, summary, content, cover, isTop, isDraft, id)
	return err
}

func (r *Repository) DeleteArticle(id int64) error {
	_, err := r.db.Exec("UPDATE article SET is_deleted = 1 WHERE id = ? AND is_deleted = 0", id)
	return err
}
