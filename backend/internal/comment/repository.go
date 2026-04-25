package comment

import (
	"database/sql"
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
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetCommentByID(id int64) (*Comment, error) {
	var comment Comment
	err := r.db.QueryRow(`SELECT id, article_id, user_id, content, created_at
    FROM comment WHERE id = ? AND is_deleted = 0`, id).Scan(&comment.ID, &comment.ArticleID, &comment.UserID, &comment.Content, &comment.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &comment, nil
}

func (r *Repository) GetCommentListByArticleID(articleID int64) ([]*CommentVO, error) {
	rows, err := r.db.Query(`
    SELECT c.id, c.article_id, c.user_id, u.username, u.avatar, c.content, c.created_at
    FROM comment c
    LEFT JOIN `+"`user`"+` u ON c.user_id = u.id
    WHERE c.article_id = ? AND c.is_deleted = 0
    ORDER BY c.created_at ASC`, articleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*CommentVO
	for rows.Next() {
		var vo CommentVO
		err := rows.Scan(&vo.ID, &vo.ArticleID, &vo.UserID, &vo.Username, &vo.Avatar, &vo.Content, &vo.CreatedAt)
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

func (r *Repository) CreateComment(comment *Comment) (int64, error) {
	result, err := r.db.Exec(`INSERT INTO comment (article_id, user_id, content) VALUES (?, ?, ?)`, comment.ArticleID, comment.UserID, comment.Content)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (r *Repository) DeleteComment(id int64) error {
	_, err := r.db.Exec(`UPDATE comment SET is_deleted = 1 WHERE id = ? AND is_deleted = 0`, id)
	return err
}
