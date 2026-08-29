package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func newTestJwtService() *JwtAuthService {
	return NewJwtAuthService("test-issuer", "test-access-secret", "test-refresh-secret", time.Hour, 24*time.Hour)
}

func TestGenerateAndParseAccessToken(t *testing.T) {
	svc := newTestJwtService()

	token, err := svc.GenerateAccessToken(42)
	if err != nil {
		t.Fatalf("签发 access token 失败: %v", err)
	}

	claims, err := svc.ParseAndVerifyAccessToken(token)
	if err != nil {
		t.Fatalf("解析 access token 失败: %v", err)
	}
	if claims.UserID != 42 || claims.TokenType != "access" || claims.Issuer != "test-issuer" {
		t.Errorf("claims 不符: userID=%d type=%s issuer=%s", claims.UserID, claims.TokenType, claims.Issuer)
	}
	if claims.JTI == "" {
		t.Error("JTI 不应为空")
	}
}

func TestRefreshTokenCannotBeUsedAsAccess(t *testing.T) {
	svc := newTestJwtService()

	refresh, err := svc.GenerateRefreshToken(1)
	if err != nil {
		t.Fatalf("签发 refresh token 失败: %v", err)
	}
	if _, err := svc.ParseAndVerifyAccessToken(refresh); err != ErrInvalidTokenType {
		t.Errorf("refresh token 用作 access 应返回 ErrInvalidTokenType，得到 %v", err)
	}
}

func TestParseWithWrongSecret(t *testing.T) {
	signer := NewJwtAuthService("test-issuer", "secret-a", "refresh-a", time.Hour, time.Hour)
	verifier := NewJwtAuthService("test-issuer", "secret-b", "refresh-a", time.Hour, time.Hour)

	token, err := signer.GenerateAccessToken(1)
	if err != nil {
		t.Fatalf("签发失败: %v", err)
	}
	if _, err := verifier.ParseAndVerifyAccessToken(token); err != ErrInvalidToken {
		t.Errorf("错误密钥验签应返回 ErrInvalidToken，得到 %v", err)
	}
}

func TestParseWithIssuerMismatch(t *testing.T) {
	signer := NewJwtAuthService("issuer-a", "secret", "refresh", time.Hour, time.Hour)
	verifier := NewJwtAuthService("issuer-b", "secret", "refresh", time.Hour, time.Hour)

	token, err := signer.GenerateAccessToken(1)
	if err != nil {
		t.Fatalf("签发失败: %v", err)
	}
	if _, err := verifier.ParseAndVerifyAccessToken(token); err != ErrInvalidToken {
		t.Errorf("issuer 不匹配应返回 ErrInvalidToken，得到 %v", err)
	}
}

func TestParseExpiredToken(t *testing.T) {
	svc := NewJwtAuthService("test-issuer", "secret", "refresh", time.Hour, time.Hour)

	// 手工签发一个已过期的 access token
	now := time.Now().Add(-2 * time.Hour)
	claims := TokenClaims{
		UserID:    1,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "test-issuer",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("secret"))
	if err != nil {
		t.Fatalf("构造过期 token 失败: %v", err)
	}

	if _, err := svc.ParseAndVerifyAccessToken(token); err != ErrInvalidToken {
		t.Errorf("过期 token 应返回 ErrInvalidToken，得到 %v", err)
	}
}

func TestCheckJwtConfig(t *testing.T) {
	cases := []struct {
		name    string
		svc     *JwtAuthService
		wantErr bool
	}{
		{"合法配置", newTestJwtService(), false},
		{"access 密钥为空", NewJwtAuthService("i", "", "r", time.Hour, time.Hour), true},
		{"refresh 密钥为空", NewJwtAuthService("i", "a", "", time.Hour, time.Hour), true},
		{"issuer 为空", NewJwtAuthService("", "a", "r", time.Hour, time.Hour), true},
		{"access TTL 为零", NewJwtAuthService("i", "a", "r", 0, time.Hour), true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.svc.CheckJwtConfig()
			if tc.wantErr && err == nil {
				t.Error("应返回错误")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("不应返回错误: %v", err)
			}
		})
	}
}
