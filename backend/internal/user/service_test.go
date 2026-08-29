package user

import (
	"context"
	"errors"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"ElainaBlog/internal/auth"
)

// fakeUserStore 实现 UserStore 接口的内存假实现
type fakeUserStore struct {
	byUsername      map[string]*User
	byEmail         map[string]*User
	created         []*User
	passwordUpdates map[int64]string
	nextID          int64
}

func newFakeUserStore() *fakeUserStore {
	return &fakeUserStore{
		byUsername:      map[string]*User{},
		byEmail:         map[string]*User{},
		passwordUpdates: map[int64]string{},
	}
}

// 编译期确认假实现满足接口
var _ UserStore = (*fakeUserStore)(nil)

func (f *fakeUserStore) GetUserByUsername(_ context.Context, username string) (*User, error) {
	u, ok := f.byUsername[username]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return u, nil
}

func (f *fakeUserStore) GetUserByEmail(_ context.Context, email string) (*User, error) {
	u, ok := f.byEmail[email]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return u, nil
}

func (f *fakeUserStore) GetUserByID(_ context.Context, id int64) (*User, error) {
	for _, u := range f.byUsername {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (f *fakeUserStore) GetUserList(context.Context) ([]*User, error) {
	var users []*User
	for _, u := range f.byUsername {
		users = append(users, u)
	}
	return users, nil
}

func (f *fakeUserStore) CreateUser(_ context.Context, u *User) (int64, error) {
	f.nextID++
	u.ID = f.nextID
	f.byUsername[u.Username] = u
	f.byEmail[u.Email] = u
	f.created = append(f.created, u)
	return u.ID, nil
}

func (f *fakeUserStore) UpdateProfile(_ context.Context, id int64, username, email, avatar string) error {
	u, err := f.GetUserByID(context.Background(), id)
	if err != nil {
		return err
	}
	delete(f.byUsername, u.Username)
	u.Username, u.Email, u.Avatar = username, email, avatar
	f.byUsername[username] = u
	return nil
}

func (f *fakeUserStore) UpdatePassword(_ context.Context, id int64, newPassword string) error {
	f.passwordUpdates[id] = newPassword
	if u, err := f.GetUserByID(context.Background(), id); err == nil {
		u.Password = newPassword
	}
	return nil
}

func (f *fakeUserStore) DeleteUser(_ context.Context, id int64) error {
	u, err := f.GetUserByID(context.Background(), id)
	if err != nil {
		return err
	}
	delete(f.byUsername, u.Username)
	delete(f.byEmail, u.Email)
	return nil
}

func newTestService(store *fakeUserStore) *Service {
	tokenMgr := auth.NewJwtAuthService("test-issuer", "test-access-secret", "test-refresh-secret", time.Hour, 24*time.Hour)
	return NewService(store, nil, tokenMgr)
}

func hashPassword(t *testing.T, plain string) string {
	t.Helper()
	hashed, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("构造密码哈希失败: %v", err)
	}
	return string(hashed)
}

func TestLogin(t *testing.T) {
	store := newFakeUserStore()
	svc := newTestService(store)
	ctx := context.Background()

	if _, err := svc.Login(ctx, LoginParams{Email: "", Password: "passw0rd123"}); !errors.Is(err, ErrInvalidLoginParams) {
		t.Errorf("空邮箱应返回 ErrInvalidLoginParams，得到 %v", err)
	}

	if _, err := svc.Login(ctx, LoginParams{Email: "not-an-email", Password: "passw0rd123"}); err == nil || errors.Is(err, ErrUserNotFound) {
		t.Errorf("非法邮箱应返回格式错误，得到 %v", err)
	}

	if _, err := svc.Login(ctx, LoginParams{Email: "nobody@example.com", Password: "passw0rd123"}); !errors.Is(err, ErrUserNotFound) {
		t.Errorf("不存在的用户应返回 ErrUserNotFound，得到 %v", err)
	}

	alice := &User{ID: 1, Username: "alice", Email: "alice@example.com", Password: hashPassword(t, "right-passw0rd")}
	store.byUsername["alice"] = alice
	store.byEmail["alice@example.com"] = alice
	if _, err := svc.Login(ctx, LoginParams{Email: "alice@example.com", Password: "wrong-passw0rd"}); !errors.Is(err, ErrPasswordMismatch) {
		t.Errorf("密码错误应返回 ErrPasswordMismatch，得到 %v", err)
	}

	result, err := svc.Login(ctx, LoginParams{Email: "alice@example.com", Password: "right-passw0rd"})
	if err != nil {
		t.Fatalf("正确密码登录失败: %v", err)
	}
	if result.UserID != 1 || result.AccessToken == "" || result.RefreshToken == "" {
		t.Errorf("登录结果不完整: %+v", result)
	}
	claims, err := svc.tokenMgr.ParseAndVerifyAccessToken(result.AccessToken)
	if err != nil || claims.UserID != 1 {
		t.Errorf("access token 解析异常: claims=%+v err=%v", claims, err)
	}
}

func TestCreateUser(t *testing.T) {
	store := newFakeUserStore()
	svc := newTestService(store)
	ctx := context.Background()

	if _, err := svc.CreateUser(ctx, CreateUserParams{Username: "", Password: "passw0rd123", Email: "a@example.com", IsAdmin: true}); !errors.Is(err, ErrInvalidParams) {
		t.Errorf("空用户名应返回 ErrInvalidParams，得到 %v", err)
	}

	if _, err := svc.CreateUser(ctx, CreateUserParams{Username: "alice", Password: "passw0rd123", Email: "bad-email", IsAdmin: true}); err == nil {
		t.Error("非法邮箱应返回错误")
	}

	if _, err := svc.CreateUser(ctx, CreateUserParams{Username: "alice", Password: "passw0rd123", Email: "a@example.com", IsAdmin: true}); err != nil {
		t.Fatalf("创建用户失败: %v", err)
	}
	if len(store.created) != 1 {
		t.Fatalf("应落库一条用户记录，得到 %d", len(store.created))
	}
	if bcrypt.CompareHashAndPassword([]byte(store.created[0].Password), []byte("passw0rd123")) != nil {
		t.Error("落库密码应为 bcrypt 哈希而非明文")
	}

	if _, err := svc.CreateUser(ctx, CreateUserParams{Username: "alice", Password: "passw0rd123", Email: "other@example.com", IsAdmin: true}); !errors.Is(err, ErrUsernameExists) {
		t.Errorf("重复用户名应返回 ErrUsernameExists，得到 %v", err)
	}
	if _, err := svc.CreateUser(ctx, CreateUserParams{Username: "bob", Password: "passw0rd123", Email: "a@example.com", IsAdmin: true}); !errors.Is(err, ErrEmailExists) {
		t.Errorf("重复邮箱应返回 ErrEmailExists，得到 %v", err)
	}
}

func TestUpdatePassword(t *testing.T) {
	store := newFakeUserStore()
	svc := newTestService(store)
	ctx := context.Background()

	store.byUsername["alice"] = &User{ID: 7, Username: "alice", Email: "alice@example.com", Password: hashPassword(t, "old-passw0rd")}

	if err := svc.UpdatePassword(ctx, 7, "wrong-passw0rd", "new-passw0rd9"); !errors.Is(err, ErrPasswordMismatch) {
		t.Errorf("旧密码错误应返回 ErrPasswordMismatch，得到 %v", err)
	}
	if err := svc.UpdatePassword(ctx, 7, "old-passw0rd", "old-passw0rd"); !errors.Is(err, ErrSamePassword) {
		t.Errorf("新旧密码相同应返回 ErrSamePassword，得到 %v", err)
	}
	if err := svc.UpdatePassword(ctx, 7, "old-passw0rd", "new-passw0rd9"); err != nil {
		t.Fatalf("改密失败: %v", err)
	}
	hashed, ok := store.passwordUpdates[7]
	if !ok {
		t.Fatal("应调用仓库层更新密码")
	}
	if bcrypt.CompareHashAndPassword([]byte(hashed), []byte("new-passw0rd9")) != nil {
		t.Error("新密码应为 bcrypt 哈希")
	}
}
