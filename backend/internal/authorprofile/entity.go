package authorprofile

// AuthorProfile 作者信息实体
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

// UpdateProfileRequest 更新作者信息请求
type UpdateProfileRequest struct {
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
