package comment

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"gorm.io/gorm"
)

// ArticleInfoProvider 获取文章信息的接口，避免直接依赖 article 模块
type ArticleInfoProvider interface {
	GetArticleAuthorInfo(ctx context.Context, id int64) (articleUserID int64, title string, err error)
}

// NotificationCreator 创建通知的接口
type NotificationCreator interface {
	CreateNotification(ctx context.Context, userID int64, nType, title, content string, targetID int64) error
}

// UserProvider 获取用户信息的接口
type UserProvider interface {
	GetUsernameByID(ctx context.Context, id int64) (string, error)
}

// CommentStore 评论数据的窄接口（消费者侧定义，参照 article 包模式），
// 便于 service 层单元测试时以假实现替换；*Repository 天然满足此接口
type CommentStore interface {
	GetCommentByID(ctx context.Context, id int64) (*Comment, error)
	GetCommentListByArticleID(ctx context.Context, articleID int64) ([]*CommentVO, error)
	GetAllCommentList(ctx context.Context) ([]*CommentVO, error)
	CreateComment(ctx context.Context, comment *Comment) (int64, error)
	DeleteComment(ctx context.Context, id int64) error
}

type Service struct {
	repo         CommentStore
	articleInfo  ArticleInfoProvider
	notifCreator NotificationCreator
	userProvider UserProvider
}

func NewService(repo CommentStore, articleInfo ArticleInfoProvider, notifCreator NotificationCreator, userProvider UserProvider) *Service {
	return &Service{repo: repo, articleInfo: articleInfo, notifCreator: notifCreator, userProvider: userProvider}
}

var (
	ErrDBNotInitialized = errors.New("数据库未初始化")
	ErrInvalidParams    = errors.New("无效的参数")
	ErrCommentNotFound  = errors.New("评论不存在")
	ErrNoPermission     = errors.New("没有权限删除此评论")
)

func (s *Service) GetCommentByID(ctx context.Context, id int64) (*Comment, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}
	if id <= 0 {
		return nil, ErrInvalidParams
	}
	c, err := s.repo.GetCommentByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCommentNotFound
		}
		return nil, err
	}
	return c, nil
}

func (s *Service) GetCommentList(ctx context.Context, articleID int64) ([]*CommentVO, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}
	if articleID <= 0 {
		return nil, ErrInvalidParams
	}
	return s.repo.GetCommentListByArticleID(ctx, articleID)
}

// GetAllCommentList 获取所有评论（管理员用）
func (s *Service) GetAllCommentList(ctx context.Context) ([]*CommentVO, error) {
	if s == nil || s.repo == nil {
		return nil, ErrDBNotInitialized
	}
	return s.repo.GetAllCommentList(ctx)
}

func (s *Service) CreateComment(ctx context.Context, params *CreateCommentParams) (int64, error) {
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
	if len(content) > 2000 {
		return 0, ErrInvalidParams
	}

	comment := &Comment{
		ArticleID: params.ArticleID,
		UserID:    params.UserID,
		Content:   content,
	}

	// 设置回复目标
	if params.ReplyToUserID != nil && *params.ReplyToUserID > 0 {
		username, err := s.userProvider.GetUsernameByID(ctx, *params.ReplyToUserID)
		if err != nil {
			return 0, ErrInvalidParams
		}
		comment.ReplyToUserID = params.ReplyToUserID
		comment.ReplyToUsername = &username
	}

	// 设置回复目标评论内容
	if params.ReplyToCommentID != nil && *params.ReplyToCommentID > 0 {
		parentComment, err := s.repo.GetCommentByID(ctx, *params.ReplyToCommentID)
		if err == nil && parentComment != nil {
			comment.ReplyToCommentID = params.ReplyToCommentID
			comment.ReplyToContent = &parentComment.Content
		}
	}

	commentID, err := s.repo.CreateComment(ctx, comment)
	if err != nil {
		return 0, err
	}

	// 异步通知（非阻塞）
	go s.notifyComment(context.Background(), params.ArticleID, params.UserID, params.ReplyToUserID, content)

	return commentID, nil
}

func (s *Service) DeleteComment(ctx context.Context, params *DeleteCommentParams, userID int64, isAdmin bool) error {
	if s == nil || s.repo == nil {
		return ErrDBNotInitialized
	}
	if params == nil || params.ID <= 0 {
		return ErrInvalidParams
	}

	// 检查评论是否存在
	comment, err := s.repo.GetCommentByID(ctx, params.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCommentNotFound
		}
		return err
	}

	// 非作者仅管理员可删除（照 article 模式，权限判定收敛在 service 层）
	if comment.UserID != userID && !isAdmin {
		return ErrNoPermission
	}

	return s.repo.DeleteComment(ctx, params.ID)
}

// notifyComment 通知相关用户有新评论
func (s *Service) notifyComment(ctx context.Context, articleID, commenterID int64, replyToUserID *int64, content string) {
	if s.notifCreator == nil {
		return
	}

	summary := content
	if len([]rune(summary)) > 50 {
		summary = string([]rune(summary)[:50]) + "..."
	}

	if replyToUserID != nil && *replyToUserID > 0 {
		// 回复某用户 → 通知被回复的人
		if *replyToUserID == commenterID {
			return // 回复自己不通知
		}
		if err := s.notifCreator.CreateNotification(ctx,
			*replyToUserID,
			"comment",
			"你的评论有了新回复",
			summary,
			articleID,
		); err != nil {
			log.Printf("创建评论回复通知失败: userID=%d articleID=%d err=%v", *replyToUserID, articleID, err)
		}
	} else {
		// 普通评论 → 通知文章作者
		if s.articleInfo == nil {
			return
		}
		articleUserID, title, err := s.articleInfo.GetArticleAuthorInfo(ctx, articleID)
		if err != nil || articleUserID == commenterID {
			return
		}
		if err := s.notifCreator.CreateNotification(ctx,
			articleUserID,
			"comment",
			fmt.Sprintf("你的文章《%s》有新评论", title),
			summary,
			articleID,
		); err != nil {
			log.Printf("创建文章评论通知失败: userID=%d articleID=%d err=%v", articleUserID, articleID, err)
		}
	}
}
