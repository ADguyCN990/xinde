package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"xinde/pkg/jwt"
	"xinde/pkg/logger"
	"xinde/pkg/response"
	"xinde/pkg/stderr"
)

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	jwtService := jwt.NewJWTService()

	return gin.HandlerFunc(func(c *gin.Context) {
		// 从请求头获取token
		token := getTokenFromHeader(c)
		if token == "" {
			response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "缺少认证token")
			c.Abort()
			return
		}

		// 验证token
		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			handleTokenError(c, err)
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UID)
		c.Set("username", claims.Username)
		c.Set("is_admin", claims.IsAdmin)
		c.Set("user_claims", claims) // 存储完整的claims

		// 继续处理请求
		c.Next()
	})
}

// AdminAuth 管理员权限中间件（需要先经过JWTAuth）
func AdminAuth() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// 检查是否已通过JWT认证
		claims, exists := c.Get("user_claims")
		if !exists {
			response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "未认证")
			c.Abort()
			return
		}

		userClaims, ok := claims.(*jwt.CustomClaims)
		if !ok {
			response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "认证信息错误")
			c.Abort()
			return
		}

		// 检查管理员权限
		if !userClaims.IsAdmin {
			response.Error(c, http.StatusForbidden, response.CodeForbidden, "需要管理员权限")
			c.Abort()
			return
		}

		c.Next()
	})
}

// OptionalAuth 可选认证中间件（token存在则验证，不存在也放行）
func OptionalAuth() gin.HandlerFunc {
	jwtService := jwt.NewJWTService()

	return gin.HandlerFunc(func(c *gin.Context) {
		token := getTokenFromHeader(c)
		if token == "" {
			// 没有token，但继续处理（游客模式）
			c.Set("is_authenticated", false)
			c.Next()
			return
		}

		// 有token，尝试验证
		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			// token无效，但不阻止请求（降级为游客模式）
			c.Set("is_authenticated", false)
			c.Next()
			return
		}

		// token有效，设置用户信息
		c.Set("is_authenticated", true)
		c.Set("user_id", claims.UID)
		c.Set("username", claims.Username)
		c.Set("is_admin", claims.IsAdmin)
		c.Set("user_claims", claims)

		c.Next()
	})
}

// getTokenFromHeader 从请求头中提取token
func getTokenFromHeader(c *gin.Context) string {
	// 支持多种token传递方式

	// 方式1: Authorization: Bearer <token>
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		// 检查是否以 "Bearer " 开头
		if strings.HasPrefix(authHeader, "Bearer ") {
			return strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

	// 方式2: 直接从 Authorization 头获取（兼容老接口）
	if authHeader != "" {
		return authHeader
	}

	// 方式3: 从查询参数获取（不推荐，但有时需要）
	return c.Query("token")
}

// handleTokenError 处理token验证错误
func handleTokenError(c *gin.Context, err error) {
	switch err.Error() {
	case stderr.ErrorTokenExpired:
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, stderr.ErrorTokenExpired)
	case stderr.ErrorTokenMalFormed:
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, stderr.ErrorTokenMalFormed)
	case stderr.ErrorTokenNotValidYet:
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, stderr.ErrorTokenNotValidYet)
	case stderr.ErrorTokenInvalid:
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, stderr.ErrorTokenInvalid)
	default:
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, stderr.ErrorTokenInvalid)
	}
	logger.Error(fmt.Sprintf("错误: %s token: %s\n", err.Error(), getTokenFromHeader(c)))
}

// GetCurrentUser 从上下文中获取当前用户信息
func GetCurrentUser(c *gin.Context) (*jwt.CustomClaims, error) {
	claims, exists := c.Get("user_claims")
	if !exists {
		return nil, fmt.Errorf("用户未认证")
	}

	userClaims, ok := claims.(*jwt.CustomClaims)
	if !ok {
		return nil, fmt.Errorf("用户信息格式错误")
	}

	return userClaims, nil
}

// GetCurrentUserID 从上下文中获取当前用户ID
func GetCurrentUserID(c *gin.Context) (uint, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, fmt.Errorf("用户未认证")
	}

	uid, ok := userID.(uint)
	if !ok {
		return 0, fmt.Errorf("用户ID格式错误")
	}

	return uid, nil
}

// IsAuthenticated 检查用户是否已认证
func IsAuthenticated(c *gin.Context) bool {
	_, exists := c.Get("user_claims")
	return exists
}

// IsAdmin 检查当前用户是否为管理员
func IsAdmin(c *gin.Context) bool {
	isAdmin, exists := c.Get("is_admin")
	if !exists {
		return false
	}

	admin, ok := isAdmin.(bool)
	return ok && admin
}
