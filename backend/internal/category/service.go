package category

import (
	"ElainaBlog/pkg/rdb"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Service struct {
	repo Repository
	rdb  rdb.RedisClient
}

func NewService(repo Repository, redis rdb.RedisClient) *Service {
	return &Service{repo: repo, rdb: redis}
}

type CreateCategoryParams struct {
	Name string
}

type UpdateCategoryParams struct {
	ID   int64
	Name string
}

var (
	ErrCategoryNotFound = errors.New("分类不存在")
	ErrCategoryExists   = errors.New("分类已存在")
	ErrDBNotInitialized = errors.New("数据库未初始化")
	ErrInvalidParams    = errors.New("无效的参数")
)

const (
	cacheKeyCategoryList = "cache:category:list"
	cacheTTLCategory     = time.Hour
)

func (s *Service) GetCategoryByID(id int64) (*CategoryVO, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}
	return s.repo.GetCategoryByID(id)
}

func (s *Service) GetCategoryByName(name string) (*CategoryVO, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}
	return s.repo.GetCategoryByName(name)
}

// GetCategoryList 获取分类列表（优先查 Redis）
func (s *Service) GetCategoryList() ([]*CategoryVO, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}

	ctx := context.Background()

	// 尝试从 Redis 读取
	if s.rdb != nil {
		val, err := s.rdb.Get(ctx, cacheKeyCategoryList).Result()
		if err == nil {
			var categories []*CategoryVO
			if json.Unmarshal([]byte(val), &categories) == nil {
				return categories, nil
			}
		}
	}

	// 缓存未命中，查库
	categories, err := s.repo.GetCategoryList()
	if err != nil {
		return nil, err
	}

	// 写入缓存
	if s.rdb != nil {
		if data, err := json.Marshal(categories); err == nil {
			s.rdb.Set(ctx, cacheKeyCategoryList, data, cacheTTLCategory)
		}
	}

	return categories, nil
}

func (s *Service) CreateCategory(params CreateCategoryParams) (*CategoryVO, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}

	categoryName := strings.TrimSpace(params.Name)
	if categoryName == "" {
		return nil, ErrInvalidParams
	}

	existing, err := s.repo.GetCategoryByName(categoryName)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, ErrCategoryExists
	}

	result, err := s.repo.CreateCategory(categoryName)
	if err == nil {
		s.invalidateCache()
	}
	return result, err
}

func (s *Service) UpdateCategory(params UpdateCategoryParams) (*CategoryVO, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}

	if params.ID <= 0 {
		return nil, ErrInvalidParams
	}

	categoryName := strings.TrimSpace(params.Name)
	if categoryName == "" {
		return nil, ErrInvalidParams
	}

	_, err := s.repo.GetCategoryByID(params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}

	existing, err := s.repo.GetCategoryByName(categoryName)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existing != nil && existing.ID != params.ID {
		return nil, ErrCategoryExists
	}

	result, err := s.repo.UpdateCategory(params.ID, categoryName)
	if err == nil {
		s.invalidateCache()
	}
	return result, err
}

func (s *Service) DeleteCategory(id int64) error {
	if s == nil || s.repo == nil {
		return ErrDBNotInitialized
	}
	if id <= 0 {
		return ErrInvalidParams
	}

	_, err := s.repo.GetCategoryByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCategoryNotFound
		}
		return err
	}

	if err := s.repo.DeleteCategory(id); err != nil {
		return err
	}
	s.invalidateCache()
	return nil
}

// ToggleTop 切换分类置顶状态（仅管理员）
func (s *Service) ToggleTop(id int64, isTop bool) error {
	if s == nil || s.repo == nil {
		return ErrDBNotInitialized
	}
	if id <= 0 {
		return ErrInvalidParams
	}

	_, err := s.repo.GetCategoryByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCategoryNotFound
		}
		return err
	}

	if err := s.repo.ToggleCategoryTop(id, isTop); err != nil {
		return err
	}
	s.invalidateCache()
	return nil
}

func (s *Service) invalidateCache() {
	if s.rdb != nil {
		s.rdb.Del(context.Background(), cacheKeyCategoryList)
	}
}
