package authorprofile

import (
	"context"

	"gorm.io/gorm"
)

// Repository 使用 GORM 存储作者信息。
type Repository struct {
	db *gorm.DB
}

// NewRepository 创建作者信息仓储实例。
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Get(ctx context.Context) (*AuthorProfile, error) {
	var p AuthorProfile
	err := r.db.WithContext(ctx).Table("author_profile").
		Where("is_deleted = 0").
		First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *Repository) Create(ctx context.Context, p *AuthorProfile) (int64, error) {
	if err := r.db.WithContext(ctx).Table("author_profile").Create(p).Error; err != nil {
		return 0, err
	}
	return p.ID, nil
}

func (r *Repository) Update(ctx context.Context, p *AuthorProfile) error {
	return r.db.WithContext(ctx).Table("author_profile").
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
