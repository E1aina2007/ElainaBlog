package friendlink

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
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
	ErrLinkNotFound     = errors.New("友链不存在")
	ErrDBNotInitialized = errors.New("数据库未初始化")
	ErrInvalidParams    = errors.New("无效的参数")
)

const (
	cacheKeyFriendLinkList = "cache:friendlink:list"
	cacheTTLFriendLink     = 24 * time.Hour
)

func normalizeURL(url string) string {
	url = strings.TrimSpace(url)
	if url != "" && !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}
	return url
}

func (s *Service) GetByID(ctx context.Context, id int64) (*FriendLinkVO, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}
	if id <= 0 {
		return nil, ErrInvalidParams
	}
	vo, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrLinkNotFound
		}
		return nil, err
	}
	return vo, nil
}

// GetList 获取友链列表（优先查 Redis）
func (s *Service) GetList(ctx context.Context) ([]*FriendLinkVO, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}

	// 尝试从 Redis 读取
	if s.rdb != nil {
		val, err := s.rdb.Get(ctx, cacheKeyFriendLinkList).Result()
		if err == nil {
			var links []*FriendLinkVO
			if json.Unmarshal([]byte(val), &links) == nil {
				return links, nil
			}
		}
	}

	// 缓存未命中，查库
	links, err := s.repo.GetList(ctx)
	if err != nil {
		return nil, err
	}

	// 写入缓存
	if s.rdb != nil {
		if data, err := json.Marshal(links); err == nil {
			s.rdb.Set(ctx, cacheKeyFriendLinkList, data, cacheTTLFriendLink)
		}
	}

	return links, nil
}

func (s *Service) Create(ctx context.Context, params CreateParams) (*FriendLinkVO, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}

	name := strings.TrimSpace(params.Name)
	url := normalizeURL(params.URL)
	if name == "" || url == "" {
		return nil, ErrInvalidParams
	}

	id, err := s.repo.Create(ctx, &FriendLink{
		Name:        name,
		URL:         url,
		Avatar:      strings.TrimSpace(params.Avatar),
		Description: strings.TrimSpace(params.Description),
		SortOrder:   params.SortOrder,
	})
	if err != nil {
		return nil, err
	}
	s.invalidateCache(ctx)
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, params UpdateParams) (*FriendLinkVO, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}
	if params.ID <= 0 {
		return nil, ErrInvalidParams
	}

	name := strings.TrimSpace(params.Name)
	url := normalizeURL(params.URL)
	if name == "" || url == "" {
		return nil, ErrInvalidParams
	}

	_, err := s.repo.GetByID(ctx, params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrLinkNotFound
		}
		return nil, err
	}

	err = s.repo.Update(ctx, &FriendLink{
		ID:          params.ID,
		Name:        name,
		URL:         url,
		Avatar:      strings.TrimSpace(params.Avatar),
		Description: strings.TrimSpace(params.Description),
		SortOrder:   params.SortOrder,
	})
	if err != nil {
		return nil, err
	}
	s.invalidateCache(ctx)
	return s.repo.GetByID(ctx, params.ID)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if s == nil || s.repo == nil {
		return ErrDBNotInitialized
	}
	if id <= 0 {
		return ErrInvalidParams
	}

	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrLinkNotFound
		}
		return err
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	s.invalidateCache(ctx)
	return nil
}

func (s *Service) invalidateCache(ctx context.Context) {
	if s.rdb != nil {
		s.rdb.Del(ctx, cacheKeyFriendLinkList)
	}
}
