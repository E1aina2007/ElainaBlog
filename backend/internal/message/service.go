package message

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

// AdminUserProvider 获取管理员用户列表的接口
type AdminUserProvider interface {
	GetAdminUserIDs() ([]int64, error)
}

// NotificationCreator 创建通知的接口
type NotificationCreator interface {
	CreateNotification(userID int64, nType, title, content string, targetID int64) error
}

type Service struct {
	repo         Repository
	adminUsers   AdminUserProvider
	notifCreator NotificationCreator
}

func NewService(repo Repository, adminUsers AdminUserProvider, notifCreator NotificationCreator) *Service {
	return &Service{repo: repo, adminUsers: adminUsers, notifCreator: notifCreator}
}

var (
	ErrDBNotInitialized = errors.New("数据库未初始化")
	ErrInvalidParams    = errors.New("无效的参数")
	ErrMessageNotFound  = errors.New("留言不存在")
)

func (s *Service) GetList(limit int) ([]*MessageVO, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}
	if limit <= 0 {
		limit = 50
	}
	return s.repo.GetList(limit)
}

func (s *Service) Create(userID int64, content string) (int64, error) {
	if s == nil || s.repo == nil {
		return 0, ErrDBNotInitialized
	}
	content = strings.TrimSpace(content)
	if content == "" {
		return 0, ErrInvalidParams
	}
	if userID <= 0 {
		return 0, ErrInvalidParams
	}
	msgID, err := s.repo.Create(&Message{
		UserID:  userID,
		Content: content,
	})
	if err != nil {
		return 0, err
	}

	// 异步通知管理员（非阻塞）
	go s.notifyAdmins(userID, content)

	return msgID, nil
}

func (s *Service) GetByID(id int64) (*Message, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}
	if id <= 0 {
		return nil, ErrInvalidParams
	}
	msg, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMessageNotFound
		}
		return nil, err
	}
	return msg, nil
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrMessageNotFound
		}
		return err
	}
	return s.repo.Delete(id)
}

// notifyAdmins 通知所有管理员有新留言
func (s *Service) notifyAdmins(userID int64, content string) {
	if s.notifCreator == nil || s.adminUsers == nil {
		return
	}

	adminIDs, err := s.adminUsers.GetAdminUserIDs()
	if err != nil {
		return
	}

	summary := content
	if len([]rune(summary)) > 50 {
		summary = string([]rune(summary)[:50]) + "..."
	}

	for _, adminID := range adminIDs {
		s.notifCreator.CreateNotification(
			adminID,
			"message",
			"你有一条新留言",
			summary,
			0,
		)
	}
}
