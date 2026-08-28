package message

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"
)

// AdminUserProvider 获取管理员用户列表的接口
type AdminUserProvider interface {
	GetAdminUserIDs(ctx context.Context) ([]int64, error)
}

// NotificationCreator 创建通知的接口
type NotificationCreator interface {
	CreateNotification(ctx context.Context, userID int64, nType, title, content string, targetID int64) error
}

type Service struct {
	repo         *Repository
	adminUsers   AdminUserProvider
	notifCreator NotificationCreator
}

func NewService(repo *Repository, adminUsers AdminUserProvider, notifCreator NotificationCreator) *Service {
	return &Service{repo: repo, adminUsers: adminUsers, notifCreator: notifCreator}
}

var (
	ErrDBNotInitialized = errors.New("数据库未初始化")
	ErrInvalidParams    = errors.New("无效的参数")
	ErrMessageNotFound  = errors.New("留言不存在")
)

func (s *Service) GetList(ctx context.Context, limit int) ([]*MessageVO, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}
	if limit <= 0 {
		limit = 50
	}
	return s.repo.GetList(ctx, limit)
}

func (s *Service) Create(ctx context.Context, userID int64, content string) (int64, error) {
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
	msgID, err := s.repo.Create(ctx, &Message{
		UserID:  userID,
		Content: content,
	})
	if err != nil {
		return 0, err
	}

	// 异步通知管理员（非阻塞）
	go s.notifyAdmins(context.Background(), userID, content)

	return msgID, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*Message, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}
	if id <= 0 {
		return nil, ErrInvalidParams
	}
	msg, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMessageNotFound
		}
		return nil, err
	}
	return msg, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if s == nil || s.repo == nil {
		return ErrDBNotInitialized
	}
	if id <= 0 {
		return ErrInvalidParams
	}
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrMessageNotFound
		}
		return err
	}
	return s.repo.Delete(ctx, id)
}

// notifyAdmins 通知所有管理员有新留言
func (s *Service) notifyAdmins(ctx context.Context, userID int64, content string) {
	if s.notifCreator == nil || s.adminUsers == nil {
		return
	}

	adminIDs, err := s.adminUsers.GetAdminUserIDs(ctx)
	if err != nil {
		return
	}

	summary := content
	if len([]rune(summary)) > 50 {
		summary = string([]rune(summary)[:50]) + "..."
	}

	for _, adminID := range adminIDs {
		s.notifCreator.CreateNotification(ctx,
			adminID,
			"message",
			"你有一条新留言",
			summary,
			0,
		)
	}
}
