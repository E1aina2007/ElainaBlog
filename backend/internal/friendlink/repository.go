package friendlink

import (
	"ElainaBlog/config/db"
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

// MySQLRepository 实现 friendlink.Repository 接口，使用 MySQL 存储。
type MySQLRepository struct {
	db db.DBTX
}

// NewRepository 创建友情链接仓储实例。
func NewRepository(db db.DBTX) *MySQLRepository {
	return &MySQLRepository{db: db}
}

func (r *MySQLRepository) GetByID(id int64) (*FriendLinkVO, error) {
	var vo FriendLinkVO
	err := r.db.QueryRow(`SELECT id, name, url, avatar, description, sort_order
		FROM friend_link WHERE id = ? AND is_deleted = 0`, id).Scan(
		&vo.ID, &vo.Name, &vo.URL, &vo.Avatar, &vo.Description, &vo.SortOrder)
	if err != nil {
		return nil, err
	}
	return &vo, nil
}

func (r *MySQLRepository) GetList() ([]*FriendLinkVO, error) {
	rows, err := r.db.Query(`SELECT id, name, url, avatar, description, sort_order
		FROM friend_link WHERE is_deleted = 0 ORDER BY sort_order DESC, id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	links := make([]*FriendLinkVO, 0)
	for rows.Next() {
		var vo FriendLinkVO
		err := rows.Scan(&vo.ID, &vo.Name, &vo.URL, &vo.Avatar, &vo.Description, &vo.SortOrder)
		if err != nil {
			return nil, err
		}
		links = append(links, &vo)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return links, nil
}

func (r *MySQLRepository) Create(link *FriendLink) (int64, error) {
	result, err := r.db.Exec(`INSERT INTO friend_link (name, url, avatar, description, sort_order)
		VALUES (?, ?, ?, ?, ?)`, link.Name, link.URL, link.Avatar, link.Description, link.SortOrder)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *MySQLRepository) Update(link *FriendLink) error {
	_, err := r.db.Exec(`UPDATE friend_link SET name = ?, url = ?, avatar = ?, description = ?, sort_order = ?
		WHERE id = ? AND is_deleted = 0`, link.Name, link.URL, link.Avatar, link.Description, link.SortOrder, link.ID)
	return err
}

func (r *MySQLRepository) Delete(id int64) error {
	_, err := r.db.Exec(`UPDATE friend_link SET is_deleted = 1 WHERE id = ? AND is_deleted = 0`, id)
	return err
}
