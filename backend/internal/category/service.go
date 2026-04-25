package category

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

func (s *Service) GetCategoryList() ([]*CategoryVO, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}

	return s.repo.GetCategoryList()
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
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if existing != nil {
		return nil, ErrCategoryExists
	}

	return s.repo.CreateCategory(categoryName)
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

	// 检查分类是否存在
	_, err := s.repo.GetCategoryByID(params.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}

	// 检查新名称是否被其他分类占用
	existing, err := s.repo.GetCategoryByName(categoryName)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if existing != nil && existing.ID != params.ID {
		return nil, ErrCategoryExists
	}

	return s.repo.UpdateCategory(params.ID, categoryName)
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
		if errors.Is(err, sql.ErrNoRows) {
			return ErrCategoryNotFound
		}
		return err
	}

	return s.repo.DeleteCategory(id)
}
