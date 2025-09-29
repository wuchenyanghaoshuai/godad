package middleware

import (
	"crypto/subtle"
	"net/http"
	"strings"

	"godad-backend/utils"

	"github.com/gin-gonic/gin"
)

var safeHTTPMethods = map[string]struct{}{
	http.MethodGet:     {},
	http.MethodHead:    {},
	http.MethodOptions: {},
	http.MethodTrace:   {},
}

// CSRFMiddleware 使用双提交策略在基于 Cookie 的会话中校验 CSRF 令牌。
func CSRFMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, ok := safeHTTPMethods[c.Request.Method]; ok {
			c.Next()
			return
		}

		// 仅在携带访问令牌 Cookie 时强制校验，避免对纯 Bearer Token 请求误报。
		if _, err := c.Cookie("access_token"); err != nil {
			c.Next()
			return
		}

		headerToken := strings.TrimSpace(c.GetHeader("X-CSRF-Token"))
		cookieToken, err := c.Cookie("csrf_token")
		if err != nil || headerToken == "" {
			utils.Error(c, utils.CodeForbidden, "缺少或无效的 CSRF 令牌")
			c.Abort()
			return
		}

		if subtle.ConstantTimeCompare([]byte(headerToken), []byte(cookieToken)) != 1 {
			utils.Error(c, utils.CodeForbidden, "CSRF 令牌验证失败")
			c.Abort()
			return
		}

		c.Next()
	}
}
