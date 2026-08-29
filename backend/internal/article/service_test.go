package article

import (
	"context"
	"errors"
	"testing"

	"gorm.io/gorm"
)

// fakeArticleStore 实现 ArticleStore 接口的内存假实现
type fakeArticleStore struct {
	byID map[int64]*ArticleVO

	listCategoryID *int64
	listSortBy     string
	listPage       int
	listPageSize   int

	searchKeyword string
	searchPage    int
	searchPageSz  int

	create createCall
	update updateCall

	toggled        map[int64]bool
	deleted        []int64
	viewIncrements []int64

	uvByID    map[int64]int64
	listRes   []*ArticleVO
	listTotal int
}

type createCall struct {
	userID     int64
	categoryID *int64
	title      string
	isTop      bool
	isDraft    bool
}

type updateCall struct {
	id      int64
	title   string
	isTop   bool
	isDraft bool
}

func (f *fakeArticleStore) GetArticleByID(_ context.Context, id int64) (*ArticleVO, error) {
	vo, ok := f.byID[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return vo, nil
}

func (f *fakeArticleStore) GetArticleByIDIncludeDraft(ctx context.Context, id int64) (*ArticleVO, error) {
	return f.GetArticleByID(ctx, id)
}

func (f *fakeArticleStore) GetArticleList(_ context.Context, categoryID *int64, sortBy string, page, pageSize int) ([]*ArticleVO, int, error) {
	f.listCategoryID, f.listSortBy, f.listPage, f.listPageSize = categoryID, sortBy, page, pageSize
	return f.listRes, f.listTotal, nil
}

func (f *fakeArticleStore) GetAdminArticleList(_ context.Context, categoryID *int64, sortBy string, page, pageSize int) ([]*ArticleVO, int, error) {
	f.listCategoryID, f.listSortBy, f.listPage, f.listPageSize = categoryID, sortBy, page, pageSize
	return f.listRes, f.listTotal, nil
}

func (f *fakeArticleStore) GetUserArticleList(_ context.Context, _ int64, categoryID *int64, sortBy string, page, pageSize int) ([]*ArticleVO, int, error) {
	f.listCategoryID, f.listSortBy, f.listPage, f.listPageSize = categoryID, sortBy, page, pageSize
	return f.listRes, f.listTotal, nil
}

func (f *fakeArticleStore) SearchArticleList(_ context.Context, keyword string, page, pageSize int) ([]*ArticleVO, int, error) {
	f.searchKeyword, f.searchPage, f.searchPageSz = keyword, page, pageSize
	return f.listRes, f.listTotal, nil
}

func (f *fakeArticleStore) CreateArticle(_ context.Context, userID int64, categoryID *int64, title, _, _, _ string, isTop, isDraft bool) (int64, error) {
	f.create = createCall{userID: userID, categoryID: categoryID, title: title, isTop: isTop, isDraft: isDraft}
	return 99, nil
}

func (f *fakeArticleStore) UpdateArticle(_ context.Context, id int64, _ *int64, title, _, _, _ string, isTop, isDraft bool) error {
	f.update = updateCall{id: id, title: title, isTop: isTop, isDraft: isDraft}
	return nil
}

func (f *fakeArticleStore) ToggleArticleTop(_ context.Context, id int64, isTop bool) error {
	f.toggled[id] = isTop
	return nil
}

func (f *fakeArticleStore) DeleteArticle(_ context.Context, id int64) error {
	f.deleted = append(f.deleted, id)
	return nil
}

func (f *fakeArticleStore) GetArticleUV(_ context.Context, id int64) (int64, error) {
	return f.uvByID[id], nil
}

func (f *fakeArticleStore) IncrementViewCountUnique(_ context.Context, id int64, _ string) error {
	f.viewIncrements = append(f.viewIncrements, id)
	return nil
}

var _ ArticleStore = (*fakeArticleStore)(nil)

type fakeCommentDeleter struct {
	deletedArticleIDs []int64
}

var _ CommentDeleter = (*fakeCommentDeleter)(nil)

func (f *fakeCommentDeleter) DeleteCommentsByArticleID(_ context.Context, articleID int64) error {
	f.deletedArticleIDs = append(f.deletedArticleIDs, articleID)
	return nil
}

func newArticleStore() *fakeArticleStore {
	return &fakeArticleStore{
		byID:    map[int64]*ArticleVO{},
		toggled: map[int64]bool{},
		uvByID:  map[int64]int64{},
	}
}

func TestGetArticleListNormalization(t *testing.T) {
	store := newArticleStore()
	svc := NewService(store, nil)
	ctx := context.Background()

	if _, err := svc.GetArticleList(ctx, &ArticleListParams{Page: 0, PageSize: 500, SortBy: "hacker"}); err != nil {
		t.Fatalf("查询失败: %v", err)
	}
	if store.listPage != 1 || store.listPageSize != 10 || store.listSortBy != "" {
		t.Errorf("非法分页与排序应归一化: page=%d pageSize=%d sortBy=%q", store.listPage, store.listPageSize, store.listSortBy)
	}

	if _, err := svc.GetArticleList(ctx, &ArticleListParams{Page: 2, PageSize: 20, SortBy: "popular"}); err != nil {
		t.Fatalf("查询失败: %v", err)
	}
	if store.listPage != 2 || store.listPageSize != 20 || store.listSortBy != "popular" {
		t.Errorf("合法参数应透传: page=%d pageSize=%d sortBy=%q", store.listPage, store.listPageSize, store.listSortBy)
	}
}

func TestGetArticleByIDNotFound(t *testing.T) {
	store := newArticleStore()
	svc := NewService(store, nil)

	if _, err := svc.GetArticleByID(context.Background(), 404); !errors.Is(err, ErrArticleNotFound) {
		t.Errorf("不存在的文章应返回 ErrArticleNotFound，得到 %v", err)
	}

	store.byID[1] = &ArticleVO{ID: 1, Title: "hello"}
	vo, err := svc.GetArticleByID(context.Background(), 1)
	if err != nil || vo.Title != "hello" {
		t.Errorf("已存在文章应原样返回: vo=%+v err=%v", vo, err)
	}
}

func TestCreateArticleValidation(t *testing.T) {
	store := newArticleStore()
	svc := NewService(store, nil)
	ctx := context.Background()

	if _, err := svc.CreateArticle(ctx, nil); !errors.Is(err, ErrInvalidParams) {
		t.Errorf("nil 参数应返回 ErrInvalidParams，得到 %v", err)
	}
	if _, err := svc.CreateArticle(ctx, &CreateArticleParams{Title: "t", Content: "c"}); !errors.Is(err, ErrInvalidParams) {
		t.Errorf("无用户应返回 ErrInvalidParams，得到 %v", err)
	}
	if _, err := svc.CreateArticle(ctx, &CreateArticleParams{UserID: 1, Title: " ", Content: "c"}); !errors.Is(err, ErrInvalidParams) {
		t.Errorf("空标题应返回 ErrInvalidParams，得到 %v", err)
	}
	// 发布（非草稿）必须选择分类
	if _, err := svc.CreateArticle(ctx, &CreateArticleParams{UserID: 1, Title: "t", Content: "c"}); !errors.Is(err, ErrInvalidParams) {
		t.Errorf("发布无分类应返回 ErrInvalidParams，得到 %v", err)
	}

	id, err := svc.CreateArticle(ctx, &CreateArticleParams{UserID: 1, Title: "  t  ", Content: "c", IsDraft: true})
	if err != nil || id != 99 {
		t.Fatalf("草稿创建失败: id=%d err=%v", id, err)
	}
	if store.create.title != "t" {
		t.Errorf("标题应去除首尾空白: %q", store.create.title)
	}
}

func TestUpdateArticlePermission(t *testing.T) {
	store := newArticleStore()
	store.byID[10] = &ArticleVO{ID: 10, UserID: 100}
	svc := NewService(store, nil)
	ctx := context.Background()

	if err := svc.UpdateArticle(ctx, nil, 100, false); !errors.Is(err, ErrInvalidParams) {
		t.Errorf("nil 参数应返回 ErrInvalidParams，得到 %v", err)
	}
	if err := svc.UpdateArticle(ctx, &UpdateArticleParams{ID: 999, Title: "t", Content: "c"}, 100, false); !errors.Is(err, ErrArticleNotFound) {
		t.Errorf("不存在的文章应返回 ErrArticleNotFound，得到 %v", err)
	}
	if err := svc.UpdateArticle(ctx, &UpdateArticleParams{ID: 10}, 100, false); !errors.Is(err, ErrInvalidParams) {
		t.Errorf("空标题应返回 ErrInvalidParams，得到 %v", err)
	}

	// 非作者且非管理员 → 无权限
	if err := svc.UpdateArticle(ctx, &UpdateArticleParams{ID: 10, Title: "t", Content: "c"}, 200, false); !errors.Is(err, ErrNoPermission) {
		t.Errorf("他人文章应返回 ErrNoPermission，得到 %v", err)
	}

	// 作者本人可更新
	if err := svc.UpdateArticle(ctx, &UpdateArticleParams{ID: 10, Title: "t", Content: "c"}, 100, false); err != nil {
		t.Fatalf("作者更新失败: %v", err)
	}
	if store.update.id != 10 || store.update.title != "t" {
		t.Errorf("更新参数不符: %+v", store.update)
	}

	// 作者置顶 → 仅管理员可置顶
	top := &UpdateArticleParams{ID: 10, Title: "t", Content: "c", IsTop: true}
	if err := svc.UpdateArticle(ctx, top, 100, false); !errors.Is(err, ErrNoPermission) {
		t.Errorf("作者置顶应返回 ErrNoPermission，得到 %v", err)
	}
	// 管理员可置顶
	if err := svc.UpdateArticle(ctx, top, 200, true); err != nil {
		t.Fatalf("管理员置顶失败: %v", err)
	}
}

func TestToggleTop(t *testing.T) {
	store := newArticleStore()
	store.byID[10] = &ArticleVO{ID: 10}
	svc := NewService(store, nil)
	ctx := context.Background()

	if err := svc.ToggleTop(ctx, 10, true, false); !errors.Is(err, ErrNoPermission) {
		t.Errorf("非管理员应返回 ErrNoPermission，得到 %v", err)
	}
	if len(store.toggled) != 0 {
		t.Error("无权限时不应触达仓库层")
	}
	if err := svc.ToggleTop(ctx, 999, true, true); !errors.Is(err, ErrArticleNotFound) {
		t.Errorf("不存在的文章应返回 ErrArticleNotFound，得到 %v", err)
	}
	if err := svc.ToggleTop(ctx, 10, true, true); err != nil {
		t.Fatalf("管理员置顶失败: %v", err)
	}
	if !store.toggled[10] {
		t.Error("置顶状态未传递到仓库层")
	}
}

func TestDeleteArticlePermission(t *testing.T) {
	store := newArticleStore()
	store.byID[10] = &ArticleVO{ID: 10, UserID: 100}
	deleter := &fakeCommentDeleter{}
	svc := NewService(store, deleter)
	ctx := context.Background()

	// 非作者且非管理员 → 无权限
	if err := svc.DeleteArticle(ctx, &DeleteArticleParams{ID: 10}, 200, false); !errors.Is(err, ErrNoPermission) {
		t.Errorf("他人文章应返回 ErrNoPermission，得到 %v", err)
	}
	if len(store.deleted) != 0 {
		t.Error("无权限时不应触达仓库层")
	}

	// 作者删除 → 级联删除评论 + 删除文章
	if err := svc.DeleteArticle(ctx, &DeleteArticleParams{ID: 10}, 100, false); err != nil {
		t.Fatalf("作者删除失败: %v", err)
	}
	if len(deleter.deletedArticleIDs) != 1 || deleter.deletedArticleIDs[0] != 10 {
		t.Errorf("应级联删除文章评论: %v", deleter.deletedArticleIDs)
	}
	if len(store.deleted) != 1 || store.deleted[0] != 10 {
		t.Errorf("应调用仓库层删除文章: %v", store.deleted)
	}

	// 管理员可删除他人文章
	store.byID[11] = &ArticleVO{ID: 11, UserID: 100}
	if err := svc.DeleteArticle(ctx, &DeleteArticleParams{ID: 11}, 200, true); err != nil {
		t.Fatalf("管理员删除失败: %v", err)
	}
}

func TestSearchArticles(t *testing.T) {
	store := newArticleStore()
	svc := NewService(store, nil)
	ctx := context.Background()

	if _, err := svc.SearchArticles(ctx, "   ", 1, 10); !errors.Is(err, ErrInvalidParams) {
		t.Errorf("空关键词应返回 ErrInvalidParams，得到 %v", err)
	}
	if _, err := svc.SearchArticles(ctx, "gorm", 0, 100); err != nil {
		t.Fatalf("搜索失败: %v", err)
	}
	if store.searchKeyword != "gorm" || store.searchPage != 1 || store.searchPageSz != 10 {
		t.Errorf("关键词与分页应归一化: keyword=%q page=%d pageSize=%d", store.searchKeyword, store.searchPage, store.searchPageSz)
	}
}

func TestIncrementViewCountAndUV(t *testing.T) {
	store := newArticleStore()
	store.uvByID[5] = 7
	svc := NewService(store, nil)
	ctx := context.Background()

	if err := svc.IncrementViewCount(ctx, 0, "1.2.3.4"); !errors.Is(err, ErrInvalidParams) {
		t.Errorf("非法文章 ID 应返回 ErrInvalidParams，得到 %v", err)
	}
	if err := svc.IncrementViewCount(ctx, 5, "1.2.3.4"); err != nil {
		t.Fatalf("浏览量递增失败: %v", err)
	}
	if len(store.viewIncrements) != 1 || store.viewIncrements[0] != 5 {
		t.Errorf("递增未传递到仓库层: %v", store.viewIncrements)
	}

	uv, err := svc.GetArticleUV(ctx, 5)
	if err != nil || uv != 7 {
		t.Errorf("UV 查询异常: uv=%d err=%v", uv, err)
	}
}
