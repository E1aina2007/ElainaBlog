package notification

import (
	"context"
	"errors"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// NotificationCreator 接口供其他模块注入，避免直接依赖 notification 模块
type NotificationCreator interface {
	CreateNotification(ctx context.Context, userID int64, nType, title, content string, targetID int64) error
}

var (
	ErrNotificationNotFound = errors.New("通知不存在")
	ErrDBNotInitialized     = errors.New("数据库未初始化")
	ErrInvalidParams        = errors.New("无效的参数")
)

// CreateNotification 创建一条通知
func (s *Service) CreateNotification(ctx context.Context, userID int64, nType, title, content string, targetID int64) error {
	if s == nil || s.repo == nil {
		return ErrDBNotInitialized
	}
	if userID <= 0 {
		return ErrInvalidParams
	}

	_, err := s.repo.Create(ctx, &Notification{
		UserID:   userID,
		Type:     nType,
		Title:    title,
		Content:  content,
		TargetID: targetID,
	})
	return err
}

// GetList 获取用户通知列表
func (s *Service) GetList(ctx context.Context, userID int64, onlyUnread bool) ([]*NotificationVO, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}
	if userID <= 0 {
		return nil, ErrInvalidParams
	}
	return s.repo.GetByUserID(ctx, userID, onlyUnread)
}

// GetUnreadCount 获取未读通知数量
func (s *Service) GetUnreadCount(ctx context.Context, userID int64) (int, error) {
	if s == nil || s.repo == nil {
		return 0, ErrDBNotInitialized
	}
	if userID <= 0 {
		return 0, ErrInvalidParams
	}
	return s.repo.GetUnreadCount(ctx, userID)
}

// MarkAsRead 标记单条通知为已读
func (s *Service) MarkAsRead(ctx context.Context, id int64, userID int64) error {
	if s == nil || s.repo == nil {
		return ErrDBNotInitialized
	}
	if id <= 0 || userID <= 0 {
		return ErrInvalidParams
	}
	return s.repo.MarkAsRead(ctx, id, userID)
}

// MarkAllAsRead 标记用户所有未读通知为已读
func (s *Service) MarkAllAsRead(ctx context.Context, userID int64) error {
	if s == nil || s.repo == nil {
		return ErrDBNotInitialized
	}
	if userID <= 0 {
		return ErrInvalidParams
	}
	return s.repo.MarkAllAsRead(ctx, userID)
}

// Delete 删除通知
func (s *Service) Delete(ctx context.Context, id int64, userID int64) error {
	if s == nil || s.repo == nil {
		return ErrDBNotInitialized
	}
	if id <= 0 || userID <= 0 {
		return ErrInvalidParams
	}
	return s.repo.Delete(ctx, id, userID)
}
