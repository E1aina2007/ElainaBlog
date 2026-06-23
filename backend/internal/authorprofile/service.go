package authorprofile

import (
	"ElainaBlog/pkg/rdb"
	"context"
	"encoding/json"
	"errors"
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

var (
	ErrDBNotInitialized = errors.New("数据库未初始化")
)

const (
	cacheKey = "cache:authorprofile"
	cacheTTL = 24 * time.Hour
)

// Get 获取作者信息（优先查 Redis）
func (s *Service) Get() (*AuthorProfile, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}

	ctx := context.Background()

	// 尝试从 Redis 读取
	if s.rdb != nil {
		val, err := s.rdb.Get(ctx, cacheKey).Result()
		if err == nil {
			var p AuthorProfile
			if json.Unmarshal([]byte(val), &p) == nil {
				return &p, nil
			}
		}
	}

	// 缓存未命中，查库
	p, err := s.repo.Get()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &AuthorProfile{
				TechStackFrontend:    "[]",
				TechStackBackend:     "[]",
				TechStackEngineering: "[]",
			}, nil
		}
		return nil, err
	}

	// 写入缓存
	if s.rdb != nil {
		if data, err := json.Marshal(p); err == nil {
			s.rdb.Set(ctx, cacheKey, data, cacheTTL)
		}
	}

	return p, nil
}

// Update 更新作者信息（写入后清除缓存）
func (s *Service) Update(p *AuthorProfile) error {
	if s == nil || s.repo == nil {
		return ErrDBNotInitialized
	}
	existing, err := s.repo.Get()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_, err = s.repo.Create(p)
			s.invalidateCache()
			return err
		}
		return err
	}
	p.ID = existing.ID
	if err := s.repo.Update(p); err != nil {
		return err
	}
	s.invalidateCache()
	return nil
}

func (s *Service) invalidateCache() {
	if s.rdb != nil {
		s.rdb.Del(context.Background(), cacheKey)
	}
}
