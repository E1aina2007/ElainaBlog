package comment

import (
	"context"
	"errors"
	"testing"

	"gorm.io/gorm"
)

// fakeCommentStore 实现 CommentStore 接口的内存假实现
type fakeCommentStore struct {
	byID            map[int64]*Comment
	articleComments map[int64][]*CommentVO
	allComments     []*CommentVO
	created         []*Comment
	deleted         []int64
	nextID          int64
}

func newFakeCommentStore() *fakeCommentStore {
	return &fakeCommentStore{
		byID:            map[int64]*Comment{},
		articleComments: map[int64][]*CommentVO{},
	}
}

var _ CommentStore = (*fakeCommentStore)(nil)

func (f *fakeCommentStore) GetCommentByID(_ context.Context, id int64) (*Comment, error) {
	c, ok := f.byID[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return c, nil
}

func (f *fakeCommentStore) GetCommentListByArticleID(_ context.Context, articleID int64) ([]*CommentVO, error) {
	return f.articleComments[articleID], nil
}

func (f *fakeCommentStore) GetAllCommentList(context.Context) ([]*CommentVO, error) {
	return f.allComments, nil
}

func (f *fakeCommentStore) CreateComment(_ context.Context, c *Comment) (int64, error) {
	f.nextID++
	c.ID = f.nextID
	f.byID[c.ID] = c
	f.created = append(f.created, c)
	return c.ID, nil
}

func (f *fakeCommentStore) DeleteComment(_ context.Context, id int64) error {
	f.deleted = append(f.deleted, id)
	delete(f.byID, id)
	return nil
}

func TestDeleteComment(t *testing.T) {
	ctx := context.Background()

	newService := func() (*Service, *fakeCommentStore) {
		store := newFakeCommentStore()
		store.byID[10] = &Comment{ID: 10, UserID: 100, Content: "hello"}
		return NewService(store, nil, nil, nil), store
	}

	if svc, _ := newService(); svc != nil {
		if err := svc.DeleteComment(ctx, nil, 100, false); !errors.Is(err, ErrInvalidParams) {
			t.Errorf("nil 参数应返回 ErrInvalidParams，得到 %v", err)
		}
	}

	svc, store := newService()
	if err := svc.DeleteComment(ctx, &DeleteCommentParams{ID: 0}, 100, false); !errors.Is(err, ErrInvalidParams) {
		t.Errorf("非法 ID 应返回 ErrInvalidParams，得到 %v", err)
	}

	if err := svc.DeleteComment(ctx, &DeleteCommentParams{ID: 999}, 100, false); !errors.Is(err, ErrCommentNotFound) {
		t.Errorf("不存在的评论应返回 ErrCommentNotFound，得到 %v", err)
	}

	// 非作者且非管理员 → 无权限，且不触达仓库层
	svc, store = newService()
	if err := svc.DeleteComment(ctx, &DeleteCommentParams{ID: 10}, 200, false); !errors.Is(err, ErrNoPermission) {
		t.Errorf("非作者删除应返回 ErrNoPermission，得到 %v", err)
	}
	if len(store.deleted) != 0 {
		t.Error("无权限时不应触达仓库层")
	}

	// 作者本人可删除
	svc, store = newService()
	if err := svc.DeleteComment(ctx, &DeleteCommentParams{ID: 10}, 100, false); err != nil {
		t.Fatalf("作者删除失败: %v", err)
	}
	if len(store.deleted) != 1 || store.deleted[0] != 10 {
		t.Errorf("删除未传递到仓库层: %v", store.deleted)
	}

	// 管理员可删除他人评论
	svc, store = newService()
	if err := svc.DeleteComment(ctx, &DeleteCommentParams{ID: 10}, 200, true); err != nil {
		t.Fatalf("管理员删除失败: %v", err)
	}
	if len(store.deleted) != 1 || store.deleted[0] != 10 {
		t.Errorf("删除未传递到仓库层: %v", store.deleted)
	}
}

func TestGetCommentByID(t *testing.T) {
	store := newFakeCommentStore()
	svc := NewService(store, nil, nil, nil)
	ctx := context.Background()

	if _, err := svc.GetCommentByID(ctx, 404); !errors.Is(err, ErrCommentNotFound) {
		t.Errorf("不存在的评论应返回 ErrCommentNotFound，得到 %v", err)
	}
	if _, err := svc.GetCommentByID(ctx, 0); !errors.Is(err, ErrInvalidParams) {
		t.Errorf("非法 ID 应返回 ErrInvalidParams，得到 %v", err)
	}

	store.byID[1] = &Comment{ID: 1, Content: "hi"}
	c, err := svc.GetCommentByID(ctx, 1)
	if err != nil || c.Content != "hi" {
		t.Errorf("已存在评论应原样返回: c=%+v err=%v", c, err)
	}
}
