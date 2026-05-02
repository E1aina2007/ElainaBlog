package user

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	Avatar    string    `json:"avatar"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// 返回给前端
type UserVO struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Avatar    string    `json:"avatar"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *User) ToVO() *UserVO {
	return &UserVO{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Avatar:    u.Avatar,
		IsAdmin:   u.IsAdmin,
		CreatedAt: u.CreatedAt,
	}
}

// AuthorInfoVO 作者公开信息（用于作者页展示）
type AuthorInfoVO struct {
	Nickname  string   `json:"nickname"`
	Avatar    string   `json:"avatar"`
	Signature string   `json:"signature"`
	Location  string   `json:"location"`
	Occupation string  `json:"occupation"`
	Hobbies   []string `json:"hobbies"`
	Email     string   `json:"email"`
	Bio       string   `json:"bio"`
	Social    struct {
		Github   string `json:"github,omitempty"`
		Blog     string `json:"blog,omitempty"`
		Bilibili string `json:"bilibili,omitempty"`
		Weibo    string `json:"weibo,omitempty"`
	} `json:"social"`
}

// AuthorStatsVO 作者统计数据
type AuthorStatsVO struct {
	ArticleCount     int `json:"article_count"`
	CommentCount     int `json:"comment_count"`
	DaysSinceCreated int `json:"days_since_created"`
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetUserByUsername(username string) (*User, error) {
	row := r.db.QueryRow("SELECT id, username, password, email, avatar, is_admin, created_at, updated_at FROM `user` WHERE username = ? AND is_deleted = 0", username)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Avatar, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserByEmail(email string) (*User, error) {
	row := r.db.QueryRow("SELECT id, username, password, email, avatar, is_admin, created_at, updated_at FROM `user` WHERE email = ? AND is_deleted = 0", email)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Avatar, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserByID(id int64) (*User, error) {
	row := r.db.QueryRow("SELECT id, username, password, email, avatar, is_admin, created_at, updated_at FROM `user` WHERE id = ? AND is_deleted = 0", id)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Avatar, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserList() ([]*User, error) {
	rows, err := r.db.Query("SELECT id, username, password, email, avatar, is_admin, created_at, updated_at FROM `user` WHERE is_deleted = 0")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Avatar, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *Repository) CreateUser(user *User) (int64, error) {
	result, err := r.db.Exec("INSERT INTO `user` (username, password, email, avatar, is_admin) VALUES (?, ?, ?, ?, ?)", user.Username, user.Password, user.Email, user.Avatar, user.IsAdmin)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *Repository) UpdateProfile(id int64, username, email, avatar string) error {
	_, err := r.db.Exec("UPDATE `user` SET username = ?, email = ?, avatar = ? WHERE id = ? AND is_deleted = 0", username, email, avatar, id)
	return err
}

func (r *Repository) UpdatePassword(id int64, newPassword string) error {
	_, err := r.db.Exec("UPDATE `user` SET password = ? WHERE id = ? AND is_deleted = 0", newPassword, id)
	return err
}

func (r *Repository) DeleteUser(id int64) error {
	_, err := r.db.Exec("UPDATE `user` SET is_deleted = 1 WHERE id = ? AND is_deleted = 0", id)
	return err
}

// GetAuthorStats 获取作者统计数据（文章数、评论数、建站天数）
func (r *Repository) GetAuthorStats(userID int64) (*AuthorStatsVO, error) {
	var stats AuthorStatsVO

	// 文章数
	err := r.db.QueryRow("SELECT COUNT(*) FROM article WHERE user_id = ? AND is_deleted = 0 AND is_draft = 0", userID).Scan(&stats.ArticleCount)
	if err != nil {
		return nil, err
	}

	// 评论数（该作者文章收到的评论）
	err = r.db.QueryRow(`
		SELECT COUNT(*) FROM comment c
		JOIN article a ON c.article_id = a.id
		WHERE a.user_id = ? AND c.is_deleted = 0`, userID).Scan(&stats.CommentCount)
	if err != nil {
		return nil, err
	}

	// 建站天数（从用户创建时间算起）
	var createdAt time.Time
	err = r.db.QueryRow("SELECT created_at FROM `user` WHERE id = ? AND is_deleted = 0", userID).Scan(&createdAt)
	if err != nil {
		return nil, err
	}
	stats.DaysSinceCreated = int(time.Since(createdAt).Hours() / 24)

	return &stats, nil
}
