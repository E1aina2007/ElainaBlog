package comment

import (
	"ElainaBlog/config/db"
	"time"
)

type Comment struct {
	ID        int64     `json:"id"`
	ArticleID int64     `json:"article_id"`
	UserID    int64     `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type CommentVO struct {
	ID        int64     `json:"id"`
	ArticleID int64     `json:"article_id"`
	UserID    int64     `json:"user_id"`
	Username  string    `json:"username"`
	Avatar    string    `json:"avatar"`
	IsAdmin   bool      `json:"is_admin"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// MySQLRepository 实现 comment.Repository 接口，使用 MySQL 存储。
type MySQLRepository struct {
	db db.DBTX
}

// NewRepository 创建评论仓储实例。
func NewRepository(db db.DBTX) *MySQLRepository {
	return &MySQLRepository{db: db}
}

func (r *MySQLRepository) GetCommentByID(id int64) (*Comment, error) {
	var comment Comment
	err := r.db.QueryRow(`SELECT id, article_id, user_id, content, created_at
    FROM comment WHERE id = ? AND is_deleted = 0`, id).Scan(&comment.ID, &comment.ArticleID, &comment.UserID, &comment.Content, &comment.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &comment, nil
}

func (r *MySQLRepository) GetCommentListByArticleID(articleID int64) ([]*CommentVO, error) {
	rows, err := r.db.Query(`
    SELECT c.id, c.article_id, c.user_id, u.username, u.avatar, u.is_admin, c.content, c.created_at
    FROM comment c
    LEFT JOIN `+"`user`"+` u ON c.user_id = u.id
    WHERE c.article_id = ? AND c.is_deleted = 0
    ORDER BY c.created_at ASC`, articleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := make([]*CommentVO, 0)
	for rows.Next() {
		var vo CommentVO
		err := rows.Scan(&vo.ID, &vo.ArticleID, &vo.UserID, &vo.Username, &vo.Avatar, &vo.IsAdmin, &vo.Content, &vo.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &vo)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *MySQLRepository) GetAllCommentList() ([]*CommentVO, error) {
	rows, err := r.db.Query(`
    SELECT c.id, c.article_id, c.user_id, u.username, u.avatar, u.is_admin, c.content, c.created_at
    FROM comment c
    LEFT JOIN ` + "`user`" + ` u ON c.user_id = u.id
    LEFT JOIN article a ON c.article_id = a.id
    WHERE c.is_deleted = 0 AND a.is_deleted = 0
    ORDER BY c.created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := make([]*CommentVO, 0)
	for rows.Next() {
		var vo CommentVO
		err := rows.Scan(&vo.ID, &vo.ArticleID, &vo.UserID, &vo.Username, &vo.Avatar, &vo.IsAdmin, &vo.Content, &vo.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &vo)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *MySQLRepository) CreateComment(comment *Comment) (int64, error) {
	result, err := r.db.Exec(`INSERT INTO comment (article_id, user_id, content) VALUES (?, ?, ?)`, comment.ArticleID, comment.UserID, comment.Content)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (r *MySQLRepository) DeleteComment(id int64) error {
	_, err := r.db.Exec(`UPDATE comment SET is_deleted = 1 WHERE id = ? AND is_deleted = 0`, id)
	return err
}

func (r *MySQLRepository) DeleteCommentsByArticleID(articleID int64) error {
	_, err := r.db.Exec(`UPDATE comment SET is_deleted = 1 WHERE article_id = ? AND is_deleted = 0`, articleID)
	return err
}
