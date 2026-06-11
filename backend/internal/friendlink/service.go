package friendlink

import (
	"database/sql"
	"errors"
	"strings"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
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

// normalizeURL 确保 URL 带有协议前缀，默认补全 https://
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrLinkNotFound
		}
		return nil, err
	}
	return vo, nil
}

func (s *Service) GetList() ([]*FriendLinkVO, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}
	return s.repo.GetList()
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

	// 检查是否存在
	_, err := s.repo.GetByID(params.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
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
		if errors.Is(err, sql.ErrNoRows) {
			return ErrLinkNotFound
		}
		return err
	}
	return s.repo.Delete(id)
}
