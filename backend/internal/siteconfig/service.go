package siteconfig

import (
	"database/sql"
	"encoding/json"
	"errors"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

var (
	ErrDBNotInitialized = errors.New("数据库未初始化")
)

// GetAllMap 返回所有配置的 key-value 映射
func (s *Service) GetAllMap() (map[string]string, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}
	configs, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	m := make(map[string]string, len(configs))
	for _, c := range configs {
		m[c.KeyName] = c.Value
	}
	return m, nil
}

// GetQuotes 返回随机语句数组
func (s *Service) GetQuotes() ([]string, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}
	c, err := s.repo.GetByKey("quotes")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []string{}, nil
		}
		return nil, err
	}
	var quotes []string
	if err := json.Unmarshal([]byte(c.Value), &quotes); err != nil {
		return []string{}, nil
	}
	return quotes, nil
}

// Upsert 批量更新配置
func (s *Service) Upsert(configs map[string]string) error {
	if s == nil || s.repo == nil {
		return ErrDBNotInitialized
	}
	for key, value := range configs {
		if err := s.repo.Upsert(key, value); err != nil {
			return err
		}
	}
	return nil
}

// Delete 删除配置
func (s *Service) Delete(key string) error {
	if s == nil || s.repo == nil {
		return ErrDBNotInitialized
	}
	return s.repo.DeleteByKey(key)
}
