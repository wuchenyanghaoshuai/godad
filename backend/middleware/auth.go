package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"godad-backend/config"
	"godad-backend/models"
	"godad-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Claims JWT声明结构
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

const (
	roleAdmin          = "admin"
	roleContentManager = "content_manager"
	roleUser           = "user"
)

// getRoleString 将角色数字转换为字符串
func getRoleString(role int8) string {
	switch role {
	case 2:
		return roleAdmin // 2-系统管理员（与现有模型保持兼容）
	case 1:
		return roleContentManager
	case 0:
		return roleUser
	default:
		return roleUser
	}
}

// GenerateToken 生成JWT令牌
func GenerateToken(user *models.User) (string, error) {
	cfg := config.GetConfig()

	// 设置过期时间
	expirationTime := time.Now().Add(time.Duration(cfg.JWT.ExpireHours) * time.Hour)

	// 创建声明
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     getRoleString(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "godad-backend",
			Subject:   strconv.Itoa(int(user.ID)),
		},
	}

	// 创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名令牌
	tokenString, err := token.SignedString([]byte(config.GetConfig().JWT.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GenerateRefreshToken 生成刷新令牌（有效期更长）
func GenerateRefreshToken(user *models.User) (string, error) {
	// 默认7天有效期
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     getRoleString(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "godad-backend",
			Subject:   strconv.Itoa(int(user.ID)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.GetConfig().JWT.Secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ParseToken 解析JWT令牌
func ParseToken(tokenString string) (*Claims, error) {
	// 解析令牌
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Header["alg"])
		}
		return []byte(config.GetConfig().JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	// 验证令牌
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrInvalidKey
}

// AuthMiddleware JWT认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 优先从 Authorization 头获取，其次从 Cookie 获取
		var tokenString string
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString = parts[1]
			}
		}
		if tokenString == "" {
			if cookie, err := c.Cookie("access_token"); err == nil {
				tokenString = cookie
			}
		}
		if tokenString == "" {
			utils.Error(c, utils.CodeUnauthorized, "缺少认证令牌")
			c.Abort()
			return
		}

		// 解析令牌
		claims, err := ParseToken(tokenString)
		if err != nil {
			utils.Error(c, utils.CodeUnauthorized, "认证令牌无效")
			c.Abort()
			return
		}

		// 检查令牌是否过期
		if claims.ExpiresAt.Time.Before(time.Now()) {
			utils.Error(c, utils.CodeUnauthorized, "认证令牌已过期")
			c.Abort()
			return
		}

		// 将用户信息存储到上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Set("claims", claims)

		c.Next()
	}
}

// OptionalAuthMiddleware 可选认证中间件（不强制要求认证）
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString = parts[1]
			}
		}
		if tokenString == "" {
			if cookie, err := c.Cookie("access_token"); err == nil {
				tokenString = cookie
			}
		}
		if tokenString == "" {
			c.Next()
			return
		}

		// 解析令牌
		claims, err := ParseToken(tokenString)
		if err != nil {
			c.Next()
			return
		}

		// 检查令牌是否过期
		if claims.ExpiresAt.Time.Before(time.Now()) {
			c.Next()
			return
		}

		// 将用户信息存储到上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Set("claims", claims)

		c.Next()
	}
}

// AdminMiddleware 管理员权限中间件
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户角色
		role, exists := c.Get("role")
		if !exists {
			utils.Error(c, utils.CodeForbidden, "权限不足")
			c.Abort()
			return
		}

		// 检查是否为管理员或内容管理员
		if role != roleAdmin && role != roleContentManager {
			utils.Error(c, utils.CodeForbidden, "需要管理员权限")
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetCurrentUserID 获取当前用户ID
func GetCurrentUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}

	if id, ok := userID.(uint); ok {
		return id, true
	}

	return 0, false
}

// GetCurrentUser 获取当前用户信息
func GetCurrentUser(c *gin.Context) (*models.User, error) {
	userID, exists := GetCurrentUserID(c)
	if !exists {
		return nil, jwt.ErrInvalidKey
	}

	// 从数据库获取用户信息
	db := config.GetDB()
	var user models.User
	if err := db.Where("id = ? AND status = ?", userID, 1).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// RefreshToken 刷新令牌
func RefreshToken(c *gin.Context) {
	// 基于 refresh_token Cookie 刷新 access_token
	refreshCookie, err := c.Cookie("refresh_token")
	if err != nil || refreshCookie == "" {
		utils.Error(c, utils.CodeUnauthorized, "缺少刷新令牌")
		return
	}
	claims, err := ParseToken(refreshCookie)
	if err != nil || claims.ExpiresAt.Time.Before(time.Now()) {
		utils.Error(c, utils.CodeUnauthorized, "刷新令牌无效或过期")
		return
	}
	// 从数据库获取用户确保有效
	db := config.GetDB()
	var user models.User
	if err := db.Where("id = ? AND status = ?", claims.UserID, 1).First(&user).Error; err != nil {
		utils.Error(c, utils.CodeUnauthorized, "用户不存在或已被禁用")
		return
	}
	access, err := GenerateToken(&user)
	if err != nil {
		utils.Error(c, utils.CodeInternalError, "生成令牌失败")
		return
	}
	// 写入新的 access_token Cookie
	csrfToken, err := SetAuthCookies(c, access, refreshCookie)
	if err != nil {
		utils.Error(c, utils.CodeInternalError, "设置认证 Cookie 失败")
		return
	}
	utils.SuccessWithMessage(c, "刷新成功", gin.H{
		"user":       user.ToResponse(),
		"csrf_token": csrfToken,
	})
}

func generateCSRFToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func resolveCookieSettings() (bool, http.SameSite) {
	cfg := config.GetConfig()
	// 默认生产环境启用 Secure
	secure := utils.GetEnvAsBool("AUTH_COOKIE_SECURE", cfg.Server.Environment == "production")

	sameSiteOption := strings.ToLower(utils.GetEnv("AUTH_COOKIE_SAMESITE", "lax"))
	var sameSite http.SameSite
	switch sameSiteOption {
	case "strict":
		sameSite = http.SameSiteStrictMode
	case "none":
		sameSite = http.SameSiteNoneMode
	default:
		sameSite = http.SameSiteLaxMode
	}

	// SameSite=None 需要 Secure=true，否则浏览器会拒绝 Cookie
	if sameSite == http.SameSiteNoneMode && !secure {
		secure = true
	}

	return secure, sameSite
}

// SetAuthCookies 设置 httpOnly Cookie (导出函数)
func SetAuthCookies(c *gin.Context, accessToken, refreshToken string) (string, error) {
	// Access Token: 短期（与 JWT 过期同步），Lax，路径根
	httpOnly := true
	secure, sameSite := resolveCookieSettings()

	// 获取JWT配置的过期时间
	cfg := config.GetConfig()
	accessTokenExpiry := time.Now().Add(time.Duration(cfg.JWT.ExpireHours) * time.Hour)
	refreshTokenExpiry := time.Now().Add(7 * 24 * time.Hour) // 7天

	csrfToken, err := generateCSRFToken()
	if err != nil {
		return "", err
	}

	// 设置 access_token
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		Expires:  accessTokenExpiry,
		MaxAge:   int(time.Duration(cfg.JWT.ExpireHours) * time.Hour / time.Second),
		HttpOnly: httpOnly,
		Secure:   secure,
		SameSite: sameSite,
	})
	// 设置 refresh_token（长效）
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		Expires:  refreshTokenExpiry,
		MaxAge:   int(7 * 24 * time.Hour / time.Second), // 7天
		HttpOnly: httpOnly,
		Secure:   secure,
		SameSite: sameSite,
	})

	// 设置 CSRF Token（非 HttpOnly，供前端读取）
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Path:     "/",
		Expires:  accessTokenExpiry,
		MaxAge:   int(time.Duration(cfg.JWT.ExpireHours) * time.Hour / time.Second),
		HttpOnly: false,
		Secure:   secure,
		SameSite: sameSite,
	})

	c.Header("X-CSRF-Token", csrfToken)
	return csrfToken, nil
}

// ClearAuthCookies 清除 Cookie (导出函数)
func ClearAuthCookies(c *gin.Context) {
	expired := time.Now().Add(-time.Hour)
	secure, sameSite := resolveCookieSettings()
	for _, name := range []string{"access_token", "refresh_token", "csrf_token"} {
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     name,
			Value:    "",
			Path:     "/",
			Expires:  expired,
			MaxAge:   -1,
			HttpOnly: name != "csrf_token",
			Secure:   secure,
			SameSite: sameSite,
		})
	}
}
