package comment

import (
	"context"

	"gorm.io/gorm"
)

const commentSelect = `c.id, c.article_id, c.user_id, c.reply_to_user_id, c.reply_to_username,
	c.reply_to_comment_id, c.reply_to_content,
	u.username, COALESCE(u.avatar,'') AS avatar, u.is_admin, c.content, c.created_at`

// Repository 使用 GORM 存储评论数据。
type Repository struct {
	db *gorm.DB
}

// NewRepository 创建评论仓储实例。
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetCommentByID(ctx context.Context, id int64) (*Comment, error) {
	var comment Comment
	err := r.db.WithContext(ctx).Table("comment").
		Where("id = ? AND is_deleted = 0", id).
		First(&comment).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *Repository) GetCommentListByArticleID(ctx context.Context, articleID int64) ([]*CommentVO, error) {
	var comments []*CommentVO
	err := r.db.WithContext(ctx).Table("comment c").
		Select(commentSelect).
		Joins("LEFT JOIN `user` u ON c.user_id = u.id").
		Where("c.article_id = ? AND c.is_deleted = 0", articleID).
		Order("c.created_at ASC").
		Scan(&comments).Error
	return comments, err
}

func (r *Repository) GetAllCommentList(ctx context.Context) ([]*CommentVO, error) {
	var comments []*CommentVO
	err := r.db.WithContext(ctx).Table("comment c").
		Select(commentSelect).
		Joins("LEFT JOIN `user` u ON c.user_id = u.id").
		Joins("LEFT JOIN article a ON c.article_id = a.id").
		Where("c.is_deleted = 0 AND a.is_deleted = 0").
		Order("c.created_at DESC").
		Scan(&comments).Error
	return comments, err
}

func (r *Repository) CreateComment(ctx context.Context, comment *Comment) (int64, error) {
	if err := r.db.WithContext(ctx).Table("comment").Create(comment).Error; err != nil {
		return 0, err
	}
	return comment.ID, nil
}

func (r *Repository) DeleteComment(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Table("comment").
		Where("id = ? AND is_deleted = 0", id).
		Update("is_deleted", 1).Error
}

func (r *Repository) DeleteCommentsByArticleID(ctx context.Context, articleID int64) error {
	return r.db.WithContext(ctx).Table("comment").
		Where("article_id = ? AND is_deleted = 0", articleID).
		Update("is_deleted", 1).Error
}
