package message

import (
	"database/sql"
	"errors"
	"strings"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
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
	return s.repo.Create(&Message{
		UserID:  userID,
		Content: content,
	})
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
		if errors.Is(err, sql.ErrNoRows) {
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
		if errors.Is(err, sql.ErrNoRows) {
			return ErrMessageNotFound
		}
		return err
	}
	return s.repo.Delete(id)
}
