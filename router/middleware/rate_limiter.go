package middleware

import (
	"net"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// 简单的滑动窗口限流：每个 IP 每分钟最多 reqPerMin 次
func RateLimiter(reqPerMin int) gin.HandlerFunc {
	type client struct {
		mu    sync.Mutex
		count int
		ts    time.Time
	}

	clients := sync.Map{}

	resetAfter := time.Minute

	return func(c *gin.Context) {
		ip := clientIP(c.Request.RemoteAddr)
		if ip == "" {
			ip = "unknown"
		}

		val, _ := clients.LoadOrStore(ip, &client{ts: time.Now()})
		cl := val.(*client)

		cl.mu.Lock()
		defer cl.mu.Unlock()

		now := time.Now()
		if now.Sub(cl.ts) > resetAfter {
			cl.count = 0
			cl.ts = now
		}
		cl.count++

		if cl.count > reqPerMin {
			c.AbortWithStatusJSON(429, gin.H{
				"error": "请求过于频繁，请稍后再试",
			})
			return
		}

		c.Next()
	}
}

func clientIP(remoteAddr string) string {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return remoteAddr
	}
	return host
}
