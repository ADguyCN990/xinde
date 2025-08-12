package jwt

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"time"
	"xinde/pkg/stderr"

	"github.com/golang-jwt/jwt/v5"
)

// JWTService handles JWT token generation and validation.
type JWTService struct {
	secretKey     []byte        // JWT 签名密钥
	tokenDuration time.Duration // Token 有效期
}

// CustomClaims defines the custom claims for our JWT.
type CustomClaims struct {
	UID      uint   `json:"uid"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

// NewJWTService creates a new JWTService.
func NewJWTService() *JWTService {
	secret := viper.GetString("jwt.secret")
	duration := viper.GetDuration("jwt.duration")
	return &JWTService{
		secretKey:     []byte(secret),
		tokenDuration: duration,
	}
}

// GenerateToken creates a new JWT token for a user.
func (s *JWTService) GenerateToken(uid uint, username string, isAdmin bool) (string, error) {
	// 创建 claims
	claims := CustomClaims{
		UID:      uid,
		Username: username,
		IsAdmin:  isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "信德刀具选型", // 签发人
		},
	}

	// 使用 HS256 签名方法创建一个新的 Token 对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥签名并获取完整的编码后的字符串 token
	signedToken, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// ValidateToken validates a JWT token string.
// It returns the custom claims if the token is valid.
func (s *JWTService) ValidateToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 确保签名方法是我们期望的 HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf(stderr.ErrorTokenInvalid)
		}
		return s.secretKey, nil
	})

	if err != nil {
		// 根据错误类型返回更具体的业务错误
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf(stderr.ErrorTokenMalFormed)
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf(stderr.ErrorTokenExpired)
		} else if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, fmt.Errorf(stderr.ErrorTokenNotValidYet)
		} else {
			return nil, fmt.Errorf(stderr.ErrorTokenInvalid)
		}
	}

	// 检查 token 是否有效，并将 claims 类型断言为 *CustomClaims
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf(stderr.ErrorTokenInvalid)
}
