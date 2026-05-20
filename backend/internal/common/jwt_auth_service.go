package common

import (
	"ElainaBlog/config"
	"ElainaBlog/pkg/rdb"
	"ElainaBlog/pkg/util"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtAuthService struct {
	AccessTokenSecret  string
	RefreshTokenSecret string
	AccessTokenTTL     time.Duration
	RefreshTokenTTL    time.Duration
	Issuer             string
}

type TokenClaims struct {
	UserID    int64  `json:"user_id"`
	TokenType string `json:"token_type"`
	JTI       string `json:"jti,omitempty"`
	jwt.RegisteredClaims
}

var (
	ErrInvalidToken     = errors.New("无效的 token")
	ErrInvalidConfig    = errors.New("jwt 配置无效")
	ErrInvalidTokenType = errors.New("无效的 token 类型")
)

// generateJTI 生成随机 JWT ID。
func generateJTI() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// BlacklistToken 将 token 的 JTI 加入 Redis 黑名单，TTL 等于 token 剩余有效期。
func BlacklistToken(jti string, ttl time.Duration) error {
	if rdb.RedisClient == nil || jti == "" {
		return nil
	}
	return rdb.RedisClient.Set(context.Background(), "token_blacklist:"+jti, "1", ttl).Err()
}

// IsTokenBlacklisted 检查 token 的 JTI 是否在黑名单中。
func IsTokenBlacklisted(jti string) bool {
	if rdb.RedisClient == nil || jti == "" {
		return false
	}
	val, err := rdb.RedisClient.Exists(context.Background(), "token_blacklist:"+jti).Result()
	return err == nil && val > 0
}

var JwtAuth *JwtAuthService

// InitJwtAuth 在配置加载完成后调用，初始化全局 JwtAuth 实例。
func InitJwtAuth() {
	JwtAuth = NewJwtAuthServiceWithConfig()
}

func NewJwtAuthService(issuer, accessTokenSecret, refreshTokenSecret string, accessTokenTTL, refreshTokenTTL time.Duration) *JwtAuthService {
	return &JwtAuthService{
		AccessTokenSecret:  accessTokenSecret,
		RefreshTokenSecret: refreshTokenSecret,
		AccessTokenTTL:     accessTokenTTL,
		RefreshTokenTTL:    refreshTokenTTL,
		Issuer:             issuer,
	}
}

func NewJwtAuthServiceWithConfig() *JwtAuthService {
	accessTTL, err := util.ParseDuration(config.GlobalConfig.Auth.AccessTokenExpiryTime)
	if err != nil {
		log.Fatalf("解析 AccessTokenExpiryTime 失败: %v", err)
	}
	refreshTTL, err := util.ParseDuration(config.GlobalConfig.Auth.RefreshTokenExpiryTime)
	if err != nil {
		log.Fatalf("解析 RefreshTokenExpiryTime 失败: %v", err)
	}

	return NewJwtAuthService(
		config.GlobalConfig.Auth.Issuer,
		config.GlobalConfig.Auth.AccessTokenSecret,
		config.GlobalConfig.Auth.RefreshTokenSecret,
		accessTTL,
		refreshTTL,
	)
}

func (s *JwtAuthService) CheckJwtConfig() error {
	if s == nil || strings.TrimSpace(s.AccessTokenSecret) == "" || strings.TrimSpace(s.RefreshTokenSecret) == "" || strings.TrimSpace(s.Issuer) == "" || s.AccessTokenTTL <= 0 || s.RefreshTokenTTL <= 0 {
		return ErrInvalidConfig
	}
	return nil
}

func (s *JwtAuthService) GenerateAccessToken(userID int64) (string, error) {
	if err := s.CheckJwtConfig(); err != nil {
		return "", err
	}

	jti, err := generateJTI()
	if err != nil {
		return "", err
	}

	now := time.Now()
	claims := TokenClaims{
		UserID:    userID,
		TokenType: "access",
		JTI:       jti,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.Issuer,
			Subject:   "user",
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.AccessTokenTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.AccessTokenSecret))
}

func (s *JwtAuthService) GenerateRefreshToken(userID int64) (string, error) {
	if err := s.CheckJwtConfig(); err != nil {
		return "", err
	}

	jti, err := generateJTI()
	if err != nil {
		return "", err
	}

	now := time.Now()
	claims := TokenClaims{
		UserID:    userID,
		TokenType: "refresh",
		JTI:       jti,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.Issuer,
			Subject:   "user",
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.RefreshTokenTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.RefreshTokenSecret))
}

func (s *JwtAuthService) ParseAndVerifyToken(tokenString string) (*TokenClaims, error) {
	if err := s.CheckJwtConfig(); err != nil {
		return nil, err
	}

	tokenString = strings.TrimSpace(tokenString)
	if tokenString == "" {
		return nil, ErrInvalidToken
	}

	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, ErrInvalidToken
		}
		// 根据 claims 中的 token_type 选择对应的密钥验签
		claims, ok := token.Claims.(*TokenClaims)
		if !ok {
			return nil, ErrInvalidToken
		}
		switch claims.TokenType {
		case "access":
			return []byte(s.AccessTokenSecret), nil
		case "refresh":
			return []byte(s.RefreshTokenSecret), nil
		default:
			return nil, ErrInvalidTokenType
		}
	})
	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	if claims.Issuer != s.Issuer {
		return nil, ErrInvalidToken
	}

	// 检查 token 是否已被吊销（加入黑名单）
	if IsTokenBlacklisted(claims.JTI) {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func (s *JwtAuthService) ParseAndVerifyAccessToken(tokenString string) (*TokenClaims, error) {
	claims, err := s.ParseAndVerifyToken(tokenString)
	if err != nil {
		return nil, err
	}
	if claims.TokenType != "access" {
		return nil, ErrInvalidTokenType
	}
	return claims, nil
}

func (s *JwtAuthService) ParseAndVerifyRefreshToken(tokenString string) (*TokenClaims, error) {
	claims, err := s.ParseAndVerifyToken(tokenString)
	if err != nil {
		return nil, err
	}
	if claims.TokenType != "refresh" {
		return nil, ErrInvalidTokenType
	}
	return claims, nil
}
