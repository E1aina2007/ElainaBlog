package site

import "errors"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

var (
	ErrDBNotInitialized = errors.New("数据库未初始化")
)

func (s *Service) GetList() (map[string]string, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}

	return s.repo.GetList()
}

var ErrInvalidParams = errors.New("无效的参数")

func (s *Service) Update(configs map[string]string) error {
	if s == nil || s.repo == nil {
		return ErrDBNotInitialized
	}
	if len(configs) == 0 {
		return ErrInvalidParams
	}
	return s.repo.Update(configs)
}
