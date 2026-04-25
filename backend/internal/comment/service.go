package comment

import (
	"database/sql"
	"errors"
	"strings"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

type CreateCommentParams struct {
	ArticleID int64
	UserID    int64
	Content   string
}

type DeleteCommentParams struct {
	ID int64
}

var (
	ErrDBNotInitialized = errors.New("数据库未初始化")
	ErrInvalidParams    = errors.New("无效的参数")
	ErrCommentNotFound  = errors.New("评论不存在")
)

func (s *Service) GetCommentByID(id int64) (*Comment, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}
	if id <= 0 {
		return nil, ErrInvalidParams
	}
	c, err := s.repo.GetCommentByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCommentNotFound
		}
		return nil, err
	}
	return c, nil
}

func (s *Service) GetCommentList(articleID int64) ([]*CommentVO, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}
	if articleID <= 0 {
		return nil, ErrInvalidParams
	}
	return s.repo.GetCommentListByArticleID(articleID)
}

func (s *Service) CreateComment(params *CreateCommentParams) (int64, error) {
	if s == nil || s.repo == nil {
		return 0, ErrDBNotInitialized
	}
	if params == nil {
		return 0, ErrInvalidParams
	}

	if params.ArticleID <= 0 || params.UserID <= 0 {
		return 0, ErrInvalidParams
	}

	content := strings.TrimSpace(params.Content)
	if content == "" {
		return 0, ErrInvalidParams
	}

	return s.repo.CreateComment(&Comment{
		ArticleID: params.ArticleID,
		UserID:    params.UserID,
		Content:   content,
	})
}

func (s *Service) DeleteComment(params *DeleteCommentParams) error {
	if s == nil || s.repo == nil {
		return ErrDBNotInitialized
	}
	if params == nil || params.ID <= 0 {
		return ErrInvalidParams
	}

	// 检查评论是否存在
	_, err := s.repo.GetCommentByID(params.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrCommentNotFound
		}
		return err
	}

	return s.repo.DeleteComment(params.ID)
}
