package authorprofile

import (
	"database/sql"
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

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Get() (*AuthorProfile, error) {
	var p AuthorProfile
	err := r.db.QueryRow(`
		SELECT id, nickname, avatar, background, signature, location, occupation, school, major,
		       email, wechat, bio, tech_stack_frontend, tech_stack_backend, tech_stack_engineering,
		       social_github, social_bilibili
		FROM author_profile WHERE is_deleted = 0 LIMIT 1
	`).Scan(
		&p.ID, &p.Nickname, &p.Avatar, &p.Background, &p.Signature, &p.Location, &p.Occupation,
		&p.School, &p.Major, &p.Email, &p.Wechat, &p.Bio, &p.TechStackFrontend,
		&p.TechStackBackend, &p.TechStackEngineering, &p.SocialGithub, &p.SocialBilibili,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *Repository) Create(p *AuthorProfile) (int64, error) {
	result, err := r.db.Exec(`
		INSERT INTO author_profile (nickname, avatar, background, signature, location, occupation,
			school, major, email, wechat, bio, tech_stack_frontend, tech_stack_backend,
			tech_stack_engineering, social_github, social_bilibili)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, p.Nickname, p.Avatar, p.Background, p.Signature, p.Location, p.Occupation,
		p.School, p.Major, p.Email, p.Wechat, p.Bio, p.TechStackFrontend,
		p.TechStackBackend, p.TechStackEngineering, p.SocialGithub, p.SocialBilibili)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *Repository) Update(p *AuthorProfile) error {
	_, err := r.db.Exec(`
		UPDATE author_profile SET
			nickname = ?, avatar = ?, background = ?, signature = ?, location = ?,
			occupation = ?, school = ?, major = ?, email = ?, wechat = ?, bio = ?,
			tech_stack_frontend = ?, tech_stack_backend = ?, tech_stack_engineering = ?,
			social_github = ?, social_bilibili = ?
		WHERE id = ? AND is_deleted = 0
	`, p.Nickname, p.Avatar, p.Background, p.Signature, p.Location, p.Occupation,
		p.School, p.Major, p.Email, p.Wechat, p.Bio, p.TechStackFrontend,
		p.TechStackBackend, p.TechStackEngineering, p.SocialGithub, p.SocialBilibili, p.ID)
	return err
}
