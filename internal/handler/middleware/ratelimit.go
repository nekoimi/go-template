package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"

	"github.com/nekoimi/go-project-template/internal/pkg/errcode"
	"github.com/nekoimi/go-project-template/internal/pkg/response"
)

// RateLimit 全局令牌桶限流中间件
func RateLimit(rps float64, burst int) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Limit(rps), burst)

	return func(c *gin.Context) {
		if !limiter.Allow() {
			response.Error(c, http.StatusTooManyRequests, errcode.TooManyReq)
			c.Abort()
			return
		}
		c.Next()
	}
}

// IPRateLimit 基于 IP 的令牌桶限流
func IPRateLimit(rps float64, burst int) gin.HandlerFunc {
	var mu sync.Mutex
	limiters := make(map[string]*rate.Limiter)

	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
		limiter, exists := limiters[ip]
		if !exists {
			limiter = rate.NewLimiter(rate.Limit(rps), burst)
			limiters[ip] = limiter
		}
		mu.Unlock()

		if !limiter.Allow() {
			response.Error(c, http.StatusTooManyRequests, errcode.TooManyReq)
			c.Abort()
			return
		}
		c.Next()
	}
}
