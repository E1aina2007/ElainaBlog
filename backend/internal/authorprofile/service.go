package authorprofile

import (
	"database/sql"
	"errors"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

var (
	ErrDBNotInitialized = errors.New("数据库未初始化")
)

// Get 获取作者信息，如果不存在则返回默认空值
func (s *Service) Get() (*AuthorProfile, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}
	p, err := s.repo.Get()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// 返回空的默认 profile，而非报错
			return &AuthorProfile{
				TechStackFrontend:    "[]",
				TechStackBackend:     "[]",
				TechStackEngineering: "[]",
			}, nil
		}
		return nil, err
	}
	return p, nil
}

// Update 更新作者信息，如果不存在则创建
func (s *Service) Update(p *AuthorProfile) error {
	if s == nil || s.repo == nil {
		return ErrDBNotInitialized
	}
	existing, err := s.repo.Get()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_, err = s.repo.Create(p)
			return err
		}
		return err
	}
	p.ID = existing.ID
	return s.repo.Update(p)
}
