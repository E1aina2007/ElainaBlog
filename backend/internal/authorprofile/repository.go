package authorprofile

import (
	"gorm.io/gorm"
)

type AuthorProfile struct {
	ID                   int64  `json:"id"`
	Nickname             string `json:"nickname"`
	Avatar               string `json:"avatar"`
	Background           string `json:"background"`
	Signature            string `json:"signature"`
	Location             string `json:"location"`
	Occupation           string `json:"occupation"`
	School               string `json:"school"`
	Major                string `json:"major"`
	Email                string `json:"email"`
	Wechat               string `json:"wechat"`
	Bio                  string `json:"bio"`
	TechStackFrontend    string `json:"tech_stack_frontend"`
	TechStackBackend     string `json:"tech_stack_backend"`
	TechStackEngineering string `json:"tech_stack_engineering"`
	SocialGithub         string `json:"social_github"`
	SocialBilibili       string `json:"social_bilibili"`
}

// MySQLRepository 实现 authorprofile.Repository 接口，使用 GORM 存储。
type MySQLRepository struct {
	db *gorm.DB
}

// NewRepository 创建作者信息仓储实例。
func NewRepository(db *gorm.DB) *MySQLRepository {
	return &MySQLRepository{db: db}
}

func (r *MySQLRepository) Get() (*AuthorProfile, error) {
	var p AuthorProfile
	err := r.db.Table("author_profile").
		Where("is_deleted = 0").
		First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *MySQLRepository) Create(p *AuthorProfile) (int64, error) {
	if err := r.db.Table("author_profile").Create(p).Error; err != nil {
		return 0, err
	}
	return p.ID, nil
}

func (r *MySQLRepository) Update(p *AuthorProfile) error {
	return r.db.Table("author_profile").
		Where("id = ? AND is_deleted = 0", p.ID).
		Updates(map[string]any{
			"nickname":              p.Nickname,
			"avatar":                p.Avatar,
			"background":            p.Background,
			"signature":             p.Signature,
			"location":              p.Location,
			"occupation":            p.Occupation,
			"school":                p.School,
			"major":                 p.Major,
			"email":                 p.Email,
			"wechat":                p.Wechat,
			"bio":                   p.Bio,
			"tech_stack_frontend":   p.TechStackFrontend,
			"tech_stack_backend":    p.TechStackBackend,
			"tech_stack_engineering": p.TechStackEngineering,
			"social_github":         p.SocialGithub,
			"social_bilibili":       p.SocialBilibili,
		}).Error
}
