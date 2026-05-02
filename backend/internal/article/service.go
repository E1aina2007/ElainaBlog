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

// ArticleListParams 文章列表查询参数
type ArticleListParams struct {
	CategoryID *int64
	Page       int
	PageSize   int
}

// ArticleListResult 文章列表返回结果
type ArticleListResult struct {
	List  []*ArticleVO `json:"list"`
	Total int          `json:"total"`
}

var (
	ErrDBNotInitialized = errors.New("数据库未初始化")
	ErrInvalidParams    = errors.New("无效的参数")
	ErrArticleNotFound  = errors.New("文章不存在")
)

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetArticleList(params *ArticleListParams) (*ArticleListResult, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}

	// 默认分页参数
	page := params.Page
	if page <= 0 {
		page = 1
	}
	pageSize := params.PageSize
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	articles, total, err := s.repo.GetArticleList(params.CategoryID, page, pageSize)
	if err != nil {
		return nil, err
	}

	return &ArticleListResult{
		List:  articles,
		Total: total,
	}, nil
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

// IncrementViewCount 增加文章浏览量
func (s *Service) IncrementViewCount(id int64) error {
	if s == nil || s.repo == nil {
		return ErrDBNotInitialized
	}
	if id <= 0 {
		return ErrInvalidParams
	}
	return s.repo.IncrementViewCount(id)
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
