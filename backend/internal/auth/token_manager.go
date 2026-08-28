// token_manager.go 定义 JWT token 管理接口，用于依赖注入和测试 mock
package auth

// TokenManager 封装 JWT token 的签发和验证操作。
// *JwtAuthService 天然满足此接口。
type TokenManager interface {
	GenerateAccessToken(userID int64) (string, error)
	GenerateRefreshToken(userID int64) (string, error)
	ParseAndVerifyAccessToken(tokenString string) (*TokenClaims, error)
	ParseAndVerifyRefreshToken(tokenString string) (*TokenClaims, error)
}
