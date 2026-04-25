package article

import (
	"database/sql"
	"errors"
	"strings"
)

type Service struct {
	repo *Repository
}

type CreateArticleParams struct {
	UserID     int64
	CategoryID *int64
	Title      string
	Summary    string
	Content    string
	Cover      string
	IsTop      bool
	IsDraft    bool
}

type UpdateArticleParams struct {
	ID         int64
	CategoryID *int64
	Title      string
	Summary    string
	Content    string
	Cover      string
	IsTop      bool
	IsDraft    bool
}

type DeleteArticleParams struct {
	ID int64
}

var (
	ErrDBNotInitialized = errors.New("数据库未初始化")
	ErrInvalidParams    = errors.New("无效的参数")
	ErrArticleNotFound  = errors.New("文章不存在")
)

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetArticleList() ([]*ArticleVO, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}
	return s.repo.GetArticleList()
}

func (s *Service) GetArticleByID(id int64) (*ArticleVO, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}
	vo, err := s.repo.GetArticleByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrArticleNotFound
		}
		return nil, err
	}
	return vo, nil
}

func (s *Service) GetArticleByUserID(userID int64) ([]*ArticleVO, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}
	return s.repo.GetArticleByUserID(userID)
}

func (s *Service) GetArticleByTitle(title string) (*ArticleVO, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}
	vo, err := s.repo.GetArticleByTitle(title)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrArticleNotFound
		}
		return nil, err
	}
	return vo, nil
}

func (s *Service) CreateArticle(params *CreateArticleParams) (int64, error) {
	if s == nil || s.repo == nil {
		return 0, ErrDBNotInitialized
	}
	if params == nil {
		return 0, ErrInvalidParams
	}

	if params.UserID <= 0 {
		return 0, ErrInvalidParams
	}

	title := strings.TrimSpace(params.Title)
	content := strings.TrimSpace(params.Content)
	if title == "" || content == "" {
		return 0, ErrInvalidParams
	}

	return s.repo.CreateArticle(
		params.UserID, params.CategoryID,
		title, strings.TrimSpace(params.Summary), content,
		strings.TrimSpace(params.Cover), params.IsTop, params.IsDraft,
	)
}

func (s *Service) UpdateArticle(params *UpdateArticleParams) error {
	if s == nil || s.repo == nil {
		return ErrDBNotInitialized
	}
	if params == nil || params.ID <= 0 {
		return ErrInvalidParams
	}

	title := strings.TrimSpace(params.Title)
	content := strings.TrimSpace(params.Content)
	if title == "" || content == "" {
		return ErrInvalidParams
	}

	// 检查文章是否存在
	_, err := s.repo.GetArticleByID(params.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrArticleNotFound
		}
		return err
	}

	return s.repo.UpdateArticle(
		params.ID, params.CategoryID,
		title, strings.TrimSpace(params.Summary), content,
		strings.TrimSpace(params.Cover), params.IsTop, params.IsDraft,
	)
}

func (s *Service) DeleteArticle(params *DeleteArticleParams) error {
	if s == nil || s.repo == nil {
		return ErrDBNotInitialized
	}
	if params == nil || params.ID <= 0 {
		return ErrInvalidParams
	}

	_, err := s.repo.GetArticleByID(params.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrArticleNotFound
		}
		return err
	}

	return s.repo.DeleteArticle(params.ID)
}
