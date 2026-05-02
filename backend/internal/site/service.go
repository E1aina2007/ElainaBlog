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

// GetAuthorInfo 从 site_config 获取作者公开信息
func (s *Service) GetAuthorInfo() (map[string]string, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}

	// 获取所有站点配置
	allConfigs, err := s.repo.GetList()
	if err != nil {
		return nil, err
	}

	// 筛选出作者相关字段
	authorFields := []string{"nickname", "signature", "location", "occupation", "hobbies", "email", "bio", "github_url", "bili_url"}
	authorInfo := make(map[string]string)
	for _, field := range authorFields {
		if val, ok := allConfigs[field]; ok {
			authorInfo[field] = val
		}
	}

	return authorInfo, nil
}

// GetDashboardStats 获取仪表盘统计数据
func (s *Service) GetDashboardStats() (*DashboardStats, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}
	return s.repo.GetDashboardStats()
}
