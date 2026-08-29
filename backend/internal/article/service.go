package article

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"
)

type CommentDeleter interface {
	DeleteCommentsByArticleID(ctx context.Context, articleID int64) error
}

// ArticleStore 文章数据的窄接口（消费者侧定义，参照 comment 包模式），
// 便于 service 层单元测试时以假实现替换；*Repository 天然满足此接口
type ArticleStore interface {
	GetArticleByID(ctx context.Context, id int64) (*ArticleVO, error)
	GetArticleByIDIncludeDraft(ctx context.Context, id int64) (*ArticleVO, error)
	GetArticleList(ctx context.Context, categoryID *int64, sortBy string, page, pageSize int) ([]*ArticleVO, int, error)
	GetAdminArticleList(ctx context.Context, categoryID *int64, sortBy string, page, pageSize int) ([]*ArticleVO, int, error)
	GetUserArticleList(ctx context.Context, userID int64, categoryID *int64, sortBy string, page, pageSize int) ([]*ArticleVO, int, error)
	SearchArticleList(ctx context.Context, keyword string, page, pageSize int) ([]*ArticleVO, int, error)
	CreateArticle(ctx context.Context, userID int64, categoryID *int64, title, summary, content, tags string, isTop, isDraft bool) (int64, error)
	UpdateArticle(ctx context.Context, id int64, categoryID *int64, title, summary, content, tags string, isTop, isDraft bool) error
	ToggleArticleTop(ctx context.Context, id int64, isTop bool) error
	DeleteArticle(ctx context.Context, id int64) error
	GetArticleUV(ctx context.Context, id int64) (int64, error)
	IncrementViewCountUnique(ctx context.Context, id int64, clientIP string) error
}

type Service struct {
	repo           ArticleStore
	commentDeleter CommentDeleter
}

var (
	ErrDBNotInitialized = errors.New("数据库未初始化")
	ErrInvalidParams    = errors.New("无效的参数")
	ErrArticleNotFound  = errors.New("文章不存在")
	ErrNoPermission     = errors.New("没有权限操作此文章")
)

func NewService(repo ArticleStore, commentDeleter CommentDeleter) *Service {
	return &Service{repo: repo, commentDeleter: commentDeleter}
}

func (s *Service) GetArticleList(ctx context.Context, params *ArticleListParams) (*ArticleListResult, error) {
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

	// 校验排序参数
	sortBy := params.SortBy
	if sortBy != "" && sortBy != "latest" && sortBy != "popular" && sortBy != "comment" {
		sortBy = ""
	}

	articles, total, err := s.repo.GetArticleList(ctx, params.CategoryID, sortBy, page, pageSize)
	if err != nil {
		return nil, err
	}

	return &ArticleListResult{
		List:  articles,
		Total: total,
	}, nil
}

func (s *Service) GetAdminArticleList(ctx context.Context, params *ArticleListParams) (*ArticleListResult, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}

	page := params.Page
	if page <= 0 {
		page = 1
	}
	pageSize := params.PageSize
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	sortBy := params.SortBy
	if sortBy != "" && sortBy != "latest" && sortBy != "popular" && sortBy != "comment" {
		sortBy = ""
	}

	articles, total, err := s.repo.GetAdminArticleList(ctx, params.CategoryID, sortBy, page, pageSize)
	if err != nil {
		return nil, err
	}

	return &ArticleListResult{
		List:  articles,
		Total: total,
	}, nil
}

func (s *Service) GetUserArticleList(ctx context.Context, userID int64, params *ArticleListParams) (*ArticleListResult, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}

	page := params.Page
	if page <= 0 {
		page = 1
	}
	pageSize := params.PageSize
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	sortBy := params.SortBy
	if sortBy != "" && sortBy != "latest" && sortBy != "popular" && sortBy != "comment" {
		sortBy = ""
	}

	articles, total, err := s.repo.GetUserArticleList(ctx, userID, params.CategoryID, sortBy, page, pageSize)
	if err != nil {
		return nil, err
	}

	return &ArticleListResult{
		List:  articles,
		Total: total,
	}, nil
}

func (s *Service) GetArticleByID(ctx context.Context, id int64) (*ArticleVO, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}
	vo, err := s.repo.GetArticleByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrArticleNotFound
		}
		return nil, err
	}
	return vo, nil
}

func (s *Service) GetArticleByIDIncludeDraft(ctx context.Context, id int64) (*ArticleVO, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}
	vo, err := s.repo.GetArticleByIDIncludeDraft(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrArticleNotFound
		}
		return nil, err
	}
	return vo, nil
}

// SearchArticles 全文搜索文章
func (s *Service) SearchArticles(ctx context.Context, keyword string, page, pageSize int) (*ArticleListResult, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}

	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return nil, ErrInvalidParams
	}
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 50 {
		pageSize = 10
	}

	articles, total, err := s.repo.SearchArticleList(ctx, keyword, page, pageSize)
	if err != nil {
		return nil, err
	}

	return &ArticleListResult{
		List:  articles,
		Total: total,
	}, nil
}

// IncrementViewCount 带 IP 去重的浏览量递增
func (s *Service) IncrementViewCount(ctx context.Context, id int64, clientIP string) error {
	if s == nil || s.repo == nil {
		return ErrDBNotInitialized
	}
	if id <= 0 {
		return ErrInvalidParams
	}

	return s.repo.IncrementViewCountUnique(ctx, id, clientIP)
}

// GetArticleUV 获取文章的独立访客数
func (s *Service) GetArticleUV(ctx context.Context, id int64) (int64, error) {
	if s == nil || s.repo == nil {
		return 0, ErrDBNotInitialized
	}
	if id <= 0 {
		return 0, ErrInvalidParams
	}
	return s.repo.GetArticleUV(ctx, id)
}

func (s *Service) CreateArticle(ctx context.Context, params *CreateArticleParams) (int64, error) {
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

	// 发布文章（非草稿）必须选择分类
	if !params.IsDraft && params.CategoryID == nil {
		return 0, ErrInvalidParams
	}

	return s.repo.CreateArticle(ctx,
		params.UserID, params.CategoryID,
		title, strings.TrimSpace(params.Summary), content, strings.TrimSpace(params.Tags),
		params.IsTop, params.IsDraft,
	)
}

func (s *Service) UpdateArticle(ctx context.Context, params *UpdateArticleParams, userID int64, isAdmin bool) error {
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

	// 检查文章是否存在（包含草稿）
	article, err := s.repo.GetArticleByIDIncludeDraft(ctx, params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrArticleNotFound
		}
		return err
	}

	// 非管理员只能操作自己的文章
	if !isAdmin && article.UserID != userID {
		return ErrNoPermission
	}

	// 只有管理员可以置顶文章
	if params.IsTop && !isAdmin {
		return ErrNoPermission
	}

	return s.repo.UpdateArticle(ctx,
		params.ID, params.CategoryID,
		title, strings.TrimSpace(params.Summary), content, strings.TrimSpace(params.Tags),
		params.IsTop, params.IsDraft,
	)
}

// ToggleTop 单独切换文章置顶状态（仅管理员）
func (s *Service) ToggleTop(ctx context.Context, id int64, isTop bool, isAdmin bool) error {
	if s == nil || s.repo == nil {
		return ErrDBNotInitialized
	}
	if id <= 0 {
		return ErrInvalidParams
	}
	if !isAdmin {
		return ErrNoPermission
	}
	// 校验文章存在
	if _, err := s.repo.GetArticleByIDIncludeDraft(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrArticleNotFound
		}
		return err
	}
	return s.repo.ToggleArticleTop(ctx, id, isTop)
}

func (s *Service) DeleteArticle(ctx context.Context, params *DeleteArticleParams, userID int64, isAdmin bool) error {
	if s == nil || s.repo == nil {
		return ErrDBNotInitialized
	}
	if params == nil || params.ID <= 0 {
		return ErrInvalidParams
	}

	// 检查文章是否存在（包含草稿）
	article, err := s.repo.GetArticleByIDIncludeDraft(ctx, params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrArticleNotFound
		}
		return err
	}

	// 非管理员只能操作自己的文章
	if !isAdmin && article.UserID != userID {
		return ErrNoPermission
	}

	// 删除文章下的所有评论
	if s.commentDeleter != nil {
		_ = s.commentDeleter.DeleteCommentsByArticleID(ctx, params.ID)
	}

	return s.repo.DeleteArticle(ctx, params.ID)
}
