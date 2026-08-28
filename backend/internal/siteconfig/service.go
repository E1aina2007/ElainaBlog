package siteconfig

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Service struct {
	repo *Repository
	rdb  *redis.Client
}

func NewService(repo *Repository, redis *redis.Client) *Service {
	return &Service{repo: repo, rdb: redis}
}

var (
	ErrDBNotInitialized = errors.New("数据库未初始化")
)

const (
	cacheKeyAll    = "cache:siteconfig:all"
	cacheKeyPrefix = "cache:siteconfig:key:"
	cacheTTL       = 24 * time.Hour
)

// GetAllMap 返回所有配置的 key-value 映射（优先查 Redis）
func (s *Service) GetAllMap(ctx context.Context) (map[string]string, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}

	// 尝试从 Redis 读取
	if s.rdb != nil {
		val, err := s.rdb.Get(ctx, cacheKeyAll).Result()
		if err == nil {
			var m map[string]string
			if json.Unmarshal([]byte(val), &m) == nil {
				return m, nil
			}
		}
	}

	// 缓存未命中，查库
	configs, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	m := make(map[string]string, len(configs))
	for _, c := range configs {
		m[c.KeyName] = c.Value
	}

	// 写入缓存
	if s.rdb != nil {
		if data, err := json.Marshal(m); err == nil {
			s.rdb.Set(ctx, cacheKeyAll, data, cacheTTL)
		}
	}

	return m, nil
}

// GetQuotes 返回随机语句数组（优先查 Redis）
func (s *Service) GetQuotes(ctx context.Context) ([]string, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}

	// 尝试从 Redis 读取
	if s.rdb != nil {
		val, err := s.rdb.Get(ctx, cacheKeyPrefix+"quotes").Result()
		if err == nil {
			var quotes []string
			if json.Unmarshal([]byte(val), &quotes) == nil {
				return quotes, nil
			}
		}
	}

	// 缓存未命中，查库
	c, err := s.repo.GetByKey(ctx, "quotes")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []string{}, nil
		}
		return nil, err
	}
	var quotes []string
	if err := json.Unmarshal([]byte(c.Value), &quotes); err != nil {
		return []string{}, nil
	}

	// 写入缓存
	if s.rdb != nil {
		if data, err := json.Marshal(quotes); err == nil {
			s.rdb.Set(ctx, cacheKeyPrefix+"quotes", data, cacheTTL)
		}
	}

	return quotes, nil
}

// Upsert 批量更新配置（写入后清除缓存）
func (s *Service) Upsert(ctx context.Context, configs map[string]string) error {
	if s == nil || s.repo == nil {
		return ErrDBNotInitialized
	}
	for key, value := range configs {
		if err := s.repo.Upsert(ctx, key, value); err != nil {
			return err
		}
	}
	s.invalidateCache(ctx)
	return nil
}

// Delete 删除配置（写入后清除缓存）
func (s *Service) Delete(ctx context.Context, key string) error {
	if s == nil || s.repo == nil {
		return ErrDBNotInitialized
	}
	if err := s.repo.DeleteByKey(ctx, key); err != nil {
		return err
	}
	s.invalidateCache(ctx)
	return nil
}

// invalidateCache 清除所有 siteconfig 缓存
func (s *Service) invalidateCache(ctx context.Context) {
	if s.rdb == nil {
		return
	}
	s.rdb.Del(ctx, cacheKeyAll, cacheKeyPrefix+"quotes")
}
