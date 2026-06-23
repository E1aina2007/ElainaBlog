package friendlink

import (
	"gorm.io/gorm"
)

// FriendLink 友情链接实体
type FriendLink struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	Avatar      string `json:"avatar"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}

// FriendLinkVO 返回给前端的视图对象
type FriendLinkVO struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	Avatar      string `json:"avatar"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}

// MySQLRepository 实现 friendlink.Repository 接口，使用 GORM 存储。
type MySQLRepository struct {
	db *gorm.DB
}

// NewRepository 创建友情链接仓储实例。
func NewRepository(db *gorm.DB) *MySQLRepository {
	return &MySQLRepository{db: db}
}

func (r *MySQLRepository) GetByID(id int64) (*FriendLinkVO, error) {
	var vo FriendLinkVO
	err := r.db.Table("friend_link").
		Select("id", "name", "url", "avatar", "description", "sort_order").
		Where("id = ? AND is_deleted = 0", id).
		Scan(&vo).Error
	if err != nil {
		return nil, err
	}
	return &vo, nil
}

func (r *MySQLRepository) GetList() ([]*FriendLinkVO, error) {
	var links []*FriendLinkVO
	err := r.db.Table("friend_link").
		Select("id", "name", "url", "avatar", "description", "sort_order").
		Where("is_deleted = 0").
		Order("sort_order DESC, id ASC").
		Scan(&links).Error
	return links, err
}

func (r *MySQLRepository) Create(link *FriendLink) (int64, error) {
	if err := r.db.Table("friend_link").Create(link).Error; err != nil {
		return 0, err
	}
	return link.ID, nil
}

func (r *MySQLRepository) Update(link *FriendLink) error {
	return r.db.Table("friend_link").
		Where("id = ? AND is_deleted = 0", link.ID).
		Updates(map[string]any{
			"name":        link.Name,
			"url":         link.URL,
			"avatar":      link.Avatar,
			"description": link.Description,
			"sort_order":  link.SortOrder,
		}).Error
}

func (r *MySQLRepository) Delete(id int64) error {
	return r.db.Table("friend_link").
		Where("id = ? AND is_deleted = 0", id).
		Update("is_deleted", 1).Error
}
