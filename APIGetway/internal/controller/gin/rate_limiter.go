package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func (api *GinAPI) RateLimiteMiddleware(limiter *rate.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		select {
		case <-c.Done():
			return
		default:
			if !limiter.Allow() {
				c.JSON(http.StatusTooManyRequests, gin.H{
					"error": "the service is temporarily unavailable",
				})
				return
			}
		}
		c.Next()
	}
}
