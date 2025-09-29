package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"godad-backend/config"
	"godad-backend/utils"

	"github.com/gin-gonic/gin"
)

func buildAllowedOrigins() ([]string, bool) {
	cfg := config.GetConfig()
	rawOrigins := utils.GetEnv("CORS_ALLOWED_ORIGINS", cfg.Server.FrontendURL)
	if rawOrigins == "" {
		return nil, false
	}

	parts := strings.Split(rawOrigins, ",")
	var origins []string
	allowWildcard := false

	for _, part := range parts {
		origin := strings.TrimSpace(part)
		if origin == "" {
			continue
		}
		if origin == "*" {
			allowWildcard = true
			continue
		}
		origins = append(origins, strings.TrimRight(origin, "/"))
	}

	return origins, allowWildcard
}

func containsOrigin(origins []string, target string) bool {
	target = strings.TrimRight(target, "/")
	for _, origin := range origins {
		if strings.EqualFold(origin, target) {
			return true
		}
	}
	return false
}

// CORSMiddleware 提供基于配置的跨域控制
func CORSMiddleware() gin.HandlerFunc {
	allowedOrigins, allowWildcard := buildAllowedOrigins()

	return func(c *gin.Context) {
		origin := strings.TrimSpace(c.GetHeader("Origin"))
		origin = strings.TrimRight(origin, "/")

		allowCredentials := false
		switch {
		case allowWildcard && origin != "":
			c.Header("Access-Control-Allow-Origin", origin)
			allowCredentials = true
		case origin != "" && containsOrigin(allowedOrigins, origin):
			c.Header("Access-Control-Allow-Origin", origin)
			allowCredentials = true
		case origin == "" && len(allowedOrigins) > 0:
			c.Header("Access-Control-Allow-Origin", allowedOrigins[0])
		}

		if allowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Vary", "Origin")
		} else {
			c.Header("Access-Control-Allow-Credentials", "false")
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-CSRF-Token, X-Requested-With, Cache-Control")
		c.Header("Access-Control-Expose-Headers", "X-CSRF-Token")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func buildDefaultCSP(cfg *config.Config) string {
	ossSources := make([]string, 0, 2)
	if cfg.OSS.CustomDomain != "" {
		ossSources = append(ossSources, cfg.OSS.CustomDomain)
	} else if cfg.OSS.BucketName != "" && cfg.OSS.Endpoint != "" {
		ossSources = append(ossSources, fmt.Sprintf("https://%s.%s", cfg.OSS.BucketName, cfg.OSS.Endpoint))
	}

	imgSources := append([]string{"'self'", "data:"}, ossSources...)
	mediaSources := append([]string{"'self'", "data:"}, ossSources...)

	return fmt.Sprintf("default-src 'self'; img-src %s; media-src %s; script-src 'self'; style-src 'self' 'unsafe-inline'; connect-src 'self'; font-src 'self' data:; frame-ancestors 'none'; form-action 'self'; base-uri 'self'",
		strings.Join(imgSources, " "),
		strings.Join(mediaSources, " "))
}

// SecurityMiddleware 设置关键安全响应头
func SecurityMiddleware() gin.HandlerFunc {
	cfg := config.GetConfig()
	defaultCSP := buildDefaultCSP(cfg)
	customCSP := utils.GetEnv("SECURITY_CSP", "")

	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		csp := defaultCSP
		if customCSP != "" {
			csp = customCSP
		}
		c.Header("Content-Security-Policy", csp)

		if cfg.Server.Environment == "production" {
			c.Header("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		}

		c.Next()
	}
}
