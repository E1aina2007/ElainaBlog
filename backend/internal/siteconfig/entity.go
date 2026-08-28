package siteconfig

// SiteConfig 站点配置实体
type SiteConfig struct {
	ID      int64  `json:"id"`
	KeyName string `json:"key_name"`
	Value   string `json:"value"`
}

// UpsertConfigsRequest 批量更新配置请求
type UpsertConfigsRequest struct {
	Configs map[string]string `json:"configs"`
}

// DeleteConfigRequest 删除配置请求
type DeleteConfigRequest struct {
	Key string `json:"key"`
}
