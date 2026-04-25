package common

import (
	"ElainaBlog/config"
	"ElainaBlog/pkg/util"
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
	jwt.RegisteredClaims
}

var (
	ErrInvalidToken     = errors.New("无效的 token")
	ErrExpiredToken     = errors.New("token 已过期")
	ErrInvalidConfig    = errors.New("jwt 配置无效")
	ErrInvalidTokenType = errors.New("无效的 token 类型")
)

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

	now := time.Now()
	claims := TokenClaims{
		UserID:    userID,
		TokenType: "access",
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

	now := time.Now()
	claims := TokenClaims{
		UserID:    userID,
		TokenType: "refresh",
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
