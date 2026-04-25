package article

import (
	"database/sql"
	"time"
)

type Article struct {
	ID         int64         `json:"id"`
	UserID     int64         `json:"user_id"`     // 作者 ID
	CategoryID sql.NullInt64 `json:"category_id"` // 所属分类 ID，可为 NULL 表示未分类
	Title      string        `json:"title"`       // 文章标题
	Summary    string        `json:"summary"`     // 文章摘要
	Content    string        `json:"content"`     // 文章正文（Markdown / HTML）
	Cover      string        `json:"cover"`       // 封面图 URL
	IsTop      bool          `json:"is_top"`      // 是否置顶
	IsDraft    bool          `json:"is_draft"`    // 是否草稿
	ViewCount  int           `json:"view_count"`  // 浏览次数
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
}

// 返回给前端
type ArticleVO struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	Username     string    `json:"username"`      // 关联查询作者名
	CategoryID   *int64    `json:"category_id"`   // nil 表示未分类
	CategoryName string    `json:"category_name"` // 关联查询分类名
	Title        string    `json:"title"`
	Summary      string    `json:"summary"`
	Content      string    `json:"content"`
	Cover        string    `json:"cover"`
	IsTop        bool      `json:"is_top"`
	IsDraft      bool      `json:"is_draft"`
	ViewCount    int       `json:"view_count"`
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
		SELECT a.id, a.user_id, u.username, a.category_id, COALESCE(c.name,''),
		       a.title, a.summary, a.content, a.cover, a.is_top, a.is_draft, a.view_count, a.created_at
		FROM article a
		LEFT JOIN `+"`user`"+` u ON a.user_id = u.id
		LEFT JOIN category c ON a.category_id = c.id AND c.is_deleted = 0
		WHERE a.id = ? AND a.is_deleted = 0`, id).Scan(
		&vo.ID, &vo.UserID, &vo.Username, &categoryID, &categoryName,
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

// GetArticleList 公开文章列表，过滤草稿
func (r *Repository) GetArticleList() ([]*ArticleVO, error) {
	rows, err := r.db.Query(`
		SELECT a.id, a.user_id, u.username, a.category_id, COALESCE(c.name,''),
		       a.title, a.summary, a.content, a.cover, a.is_top, a.is_draft, a.view_count, a.created_at
		FROM article a
		LEFT JOIN ` + "`user`" + ` u ON a.user_id = u.id
		LEFT JOIN category c ON a.category_id = c.id AND c.is_deleted = 0
		WHERE a.is_deleted = 0 AND a.is_draft = 0
		ORDER BY a.is_top DESC, a.created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []*ArticleVO
	for rows.Next() {
		var vo ArticleVO
		var categoryID sql.NullInt64
		var categoryName string
		err := rows.Scan(&vo.ID, &vo.UserID, &vo.Username, &categoryID, &categoryName,
			&vo.Title, &vo.Summary, &vo.Content, &vo.Cover, &vo.IsTop, &vo.IsDraft, &vo.ViewCount, &vo.CreatedAt)
		if err != nil {
			return nil, err
		}
		if categoryID.Valid {
			vo.CategoryID = &categoryID.Int64
			vo.CategoryName = categoryName
		}
		articles = append(articles, &vo)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return articles, nil
}

func (r *Repository) GetArticleByUserID(userID int64) ([]*ArticleVO, error) {
	rows, err := r.db.Query(`
		SELECT a.id, a.user_id, u.username, a.category_id, COALESCE(c.name,''),
		       a.title, a.summary, a.content, a.cover, a.is_top, a.is_draft, a.view_count, a.created_at
		FROM article a
		LEFT JOIN `+"`user`"+` u ON a.user_id = u.id
		LEFT JOIN category c ON a.category_id = c.id AND c.is_deleted = 0
		WHERE a.user_id = ? AND a.is_deleted = 0
		ORDER BY a.created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []*ArticleVO
	for rows.Next() {
		var vo ArticleVO
		var categoryID sql.NullInt64
		var categoryName string
		err := rows.Scan(&vo.ID, &vo.UserID, &vo.Username, &categoryID, &categoryName,
			&vo.Title, &vo.Summary, &vo.Content, &vo.Cover, &vo.IsTop, &vo.IsDraft, &vo.ViewCount, &vo.CreatedAt)
		if err != nil {
			return nil, err
		}
		if categoryID.Valid {
			vo.CategoryID = &categoryID.Int64
			vo.CategoryName = categoryName
		}
		articles = append(articles, &vo)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return articles, nil
}

func (r *Repository) GetArticleByTitle(title string) (*ArticleVO, error) {
	var vo ArticleVO
	var categoryID sql.NullInt64
	var categoryName string
	err := r.db.QueryRow(`
		SELECT a.id, a.user_id, u.username, a.category_id, COALESCE(c.name,''),
		       a.title, a.summary, a.content, a.cover, a.is_top, a.is_draft, a.view_count, a.created_at
		FROM article a
		LEFT JOIN `+"`user`"+` u ON a.user_id = u.id
		LEFT JOIN category c ON a.category_id = c.id AND c.is_deleted = 0
		WHERE a.title = ? AND a.is_deleted = 0`, title).Scan(
		&vo.ID, &vo.UserID, &vo.Username, &categoryID, &categoryName,
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
