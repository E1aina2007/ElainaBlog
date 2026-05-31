package comment

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// ArticleInfoProvider 获取文章信息的接口，避免直接依赖 article 模块
type ArticleInfoProvider interface {
	GetArticleAuthorInfo(id int64) (articleUserID int64, title string, err error)
}

// NotificationCreator 创建通知的接口
type NotificationCreator interface {
	CreateNotification(userID int64, nType, title, content string, targetID int64) error
}

type Service struct {
	repo           Repository
	articleInfo    ArticleInfoProvider
	notifCreator   NotificationCreator
}

func NewService(repo Repository, articleInfo ArticleInfoProvider, notifCreator NotificationCreator) *Service {
	return &Service{repo: repo, articleInfo: articleInfo, notifCreator: notifCreator}
}

type CreateCommentParams struct {
	ArticleID int64
	UserID    int64
	Content   string
}

type DeleteCommentParams struct {
	ID int64
}

var (
	ErrDBNotInitialized = errors.New("数据库未初始化")
	ErrInvalidParams    = errors.New("无效的参数")
	ErrCommentNotFound  = errors.New("评论不存在")
)

func (s *Service) GetCommentByID(id int64) (*Comment, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}
	if id <= 0 {
		return nil, ErrInvalidParams
	}
	c, err := s.repo.GetCommentByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCommentNotFound
		}
		return nil, err
	}
	return c, nil
}

func (s *Service) GetCommentList(articleID int64) ([]*CommentVO, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}
	if articleID <= 0 {
		return nil, ErrInvalidParams
	}
	return s.repo.GetCommentListByArticleID(articleID)
}

// GetAllCommentList 获取所有评论（管理员用）
func (s *Service) GetAllCommentList() ([]*CommentVO, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}
	return s.repo.GetAllCommentList()
}

func (s *Service) CreateComment(params *CreateCommentParams) (int64, error) {
	if s == nil || s.repo == nil {
		return 0, ErrDBNotInitialized
	}
	if params == nil {
		return 0, ErrInvalidParams
	}

	if params.ArticleID <= 0 || params.UserID <= 0 {
		return 0, ErrInvalidParams
	}

	content := strings.TrimSpace(params.Content)
	if content == "" {
		return 0, ErrInvalidParams
	}

	commentID, err := s.repo.CreateComment(&Comment{
		ArticleID: params.ArticleID,
		UserID:    params.UserID,
		Content:   content,
	})
	if err != nil {
		return 0, err
	}

	// 异步通知文章作者（非阻塞）
	go s.notifyArticleAuthor(params.ArticleID, params.UserID, content)

	return commentID, nil
}

func (s *Service) DeleteComment(params *DeleteCommentParams) error {
	if s == nil || s.repo == nil {
		return ErrDBNotInitialized
	}
	if params == nil || params.ID <= 0 {
		return ErrInvalidParams
	}

	// 检查评论是否存在
	_, err := s.repo.GetCommentByID(params.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrCommentNotFound
		}
		return err
	}

	return s.repo.DeleteComment(params.ID)
}

// notifyArticleAuthor 通知文章作者有新评论
func (s *Service) notifyArticleAuthor(articleID, commentUserID int64, commentContent string) {
	if s.notifCreator == nil || s.articleInfo == nil {
		return
	}

	articleUserID, title, err := s.articleInfo.GetArticleAuthorInfo(articleID)
	if err != nil {
		return
	}

	// 不通知自己
	if articleUserID == commentUserID {
		return
	}

	summary := commentContent
	if len([]rune(summary)) > 50 {
		summary = string([]rune(summary)[:50]) + "..."
	}

	s.notifCreator.CreateNotification(
		articleUserID,
		"comment",
		fmt.Sprintf("你的文章《%s》有新评论", title),
		summary,
		articleID,
	)
}
