package friendlink

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

type CreateParams struct {
	Name        string
	URL         string
	Avatar      string
	Description string
	SortOrder   int
}

type UpdateParams struct {
	ID          int64
	Name        string
	URL         string
	Avatar      string
	Description string
	SortOrder   int
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

func (s *Service) GetByID(id int64) (*FriendLinkVO, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}
	if id <= 0 {
		return nil, ErrInvalidParams
	}
	vo, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrLinkNotFound
		}
		return nil, err
	}
	return vo, nil
}

// GetList 获取友链列表（优先查 Redis）
func (s *Service) GetList() ([]*FriendLinkVO, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}

	ctx := context.Background()

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
	links, err := s.repo.GetList()
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

func (s *Service) Create(params CreateParams) (*FriendLinkVO, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}

	name := strings.TrimSpace(params.Name)
	url := normalizeURL(params.URL)
	if name == "" || url == "" {
		return nil, ErrInvalidParams
	}

	id, err := s.repo.Create(&FriendLink{
		Name:        name,
		URL:         url,
		Avatar:      strings.TrimSpace(params.Avatar),
		Description: strings.TrimSpace(params.Description),
		SortOrder:   params.SortOrder,
	})
	if err != nil {
		return nil, err
	}
	s.invalidateCache()
	return s.repo.GetByID(id)
}

func (s *Service) Update(params UpdateParams) (*FriendLinkVO, error) {
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

	_, err := s.repo.GetByID(params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrLinkNotFound
		}
		return nil, err
	}

	err = s.repo.Update(&FriendLink{
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
	s.invalidateCache()
	return s.repo.GetByID(params.ID)
}

func (s *Service) Delete(id int64) error {
	if s == nil || s.repo == nil {
		return ErrDBNotInitialized
	}
	if id <= 0 {
		return ErrInvalidParams
	}

	_, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrLinkNotFound
		}
		return err
	}
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	s.invalidateCache()
	return nil
}

func (s *Service) invalidateCache() {
	if s.rdb != nil {
		s.rdb.Del(context.Background(), cacheKeyFriendLinkList)
	}
}
